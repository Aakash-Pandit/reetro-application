package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Aakash-Pandit/reetro-golang/common"
	"github.com/Aakash-Pandit/reetro-golang/core"
	"github.com/Aakash-Pandit/reetro-golang/middlewares"
	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/Aakash-Pandit/reetro-golang/storages"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Store       storages.Storage
	User        RequestUser
	RedisClient storages.RedisStoreInterface
	Email       common.EmailInterface
}

func NewUserService(store storages.Storage, user RequestUser, redisClient storages.RedisStoreInterface, emailInterface common.EmailInterface) *UserService {
	return &UserService{Store: store, User: user, RedisClient: redisClient, Email: emailInterface}
}

func (u *UserService) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) error {
	limit, offset := core.Pagination(r)
	users, err := u.Store.GetAllUsers(limit, offset)
	if err != nil {
		log.Println("Error in fetching the users", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	user := u.User.GetRequestUser(r)
	if user.UserType != models.SuperAdmin {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusUnauthorized,
			Data:   &core.APIError{Detail: "Unauthorized Access"},
		})
	}

	return core.ListAPIResponse(w, &core.ListAPI{
		Status: http.StatusOK,
		Result: &core.ListAPIResponseBody{
			Count:  len(users),
			Result: users,
		},
	})
}

func (u *UserService) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		log.Println("Error in parsing the id", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: err.Error()},
		})
	}

	var cacheUser models.CreateUserResponse
	user, err := u.RedisClient.Get(id.String(), cacheUser)
	if user != nil {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusOK,
			Data:   user,
		})
	}

	if err != nil {
		log.Println("Error in fetching the user from cache", err)
	}

	user, err = u.Store.GetUserById(id)
	if err != nil {
		log.Println("Error in fetching the user", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "User not found"},
		})
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (u *UserService) CreateUserHandler(w http.ResponseWriter, r *http.Request) error {
	var userRequest models.CreateUserRequest

	json.NewDecoder(r.Body).Decode(&userRequest)
	err := models.ValidateStruct(&userRequest)
	if err != nil {
		log.Println("Error in validating the user struct", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	defer r.Body.Close()

	user, bcrypt_err := models.NewUser(&userRequest)
	if bcrypt_err != nil {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   bcrypt_err,
		})
	}

	newUser, store_error := u.Store.CreateUser(*user)
	if store_error != nil {
		msg := common.AnyToAnyStructField(store_error, &core.DatabaseError{})

		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   msg,
		})
	}

	userResponse := models.CreateUserResponse{}
	common.AnyToAnyStructField(newUser, &userResponse)

	redisErr := u.RedisClient.Set(newUser.Id.String(), userResponse)
	if redisErr != nil {
		log.Println("Error in setting the user in redis", redisErr)
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusCreated,
		Data:   userResponse,
	})
}

func (u *UserService) UpdateUserHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *UserService) DeleteUserHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		log.Println("Error in parsing the id", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: err.Error()},
		})
	}

	err = u.Store.DeleteUser(id)
	if err != nil {
		log.Println("Error in fetching the user:", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "User not found"},
		})
	}

	redisErr := u.RedisClient.Del(id.String())
	if redisErr != nil {
		log.Println("Error in setting the user in redis", redisErr)
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   map[string]string{"detail": "User deleted successfully"},
	})
}

func (u *UserService) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	var loginRequest models.LoginRequest

	json.NewDecoder(r.Body).Decode(&loginRequest)
	err := models.ValidateStruct(&loginRequest)
	if err != nil {
		log.Println("Error in validating the user struct", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	defer r.Body.Close()

	user, dbError := u.Store.VerifyUserByUsername(loginRequest.Username)
	if dbError != nil {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Invalid Username"},
		})
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if passwordError != nil {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Invalid Password"},
		})
	}

	token, tokenErr := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)
	if tokenErr != nil {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Error While Generating Token"},
		})
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   &core.TokenResponse{Token: token},
	})
}

func (u *UserService) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) error {
	var payload ForgotPasswordPayload

	json.NewDecoder(r.Body).Decode(&payload)
	err := models.ValidateStruct(payload)
	if err != nil {
		log.Println("Error in validating the user struct", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	defer r.Body.Close()

	user, dbError := u.Store.VerifyUserByEmail(payload.Email)
	if dbError != nil {
		log.Println("Error in fetching the user", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "User not found"},
		})
	}

	newPassword := common.GenerateRandomPassword()
	password, bcrypt_err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if bcrypt_err != nil {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Error While Generating Password, Try again later"},
		})
	}

	user.Password = string(password)

	dbError = u.Store.SavePassword(*user)
	if dbError != nil {
		log.Println("Error while updating password", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Error while Updating Password"},
		})
	}

	emailErr := u.Email.SendEmailForPasswordReset(user.Email, "Reset Password", newPassword)
	if emailErr != nil {
		log.Println("Error while sending email", emailErr)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Error while Updating Password"},
		})
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   &APISuccessResponse{Detail: "Your Password is Successfully Updated, Please check your email"},
	})
}

func (u *UserService) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) error {
	var payload ResetPasswordPayload

	json.NewDecoder(r.Body).Decode(&payload)
	err := models.ValidateStruct(payload)
	if err != nil {
		log.Println("Error in validating the user struct", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	defer r.Body.Close()

	user, dbErr := u.Store.VerifyUserByUsernamePassword(payload.Username, payload.OldPassword)
	if dbErr != nil {
		log.Println("Error in fetching the user", dbErr)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Error while Fetching User"},
		})
	}

	bcryptPassword, bcryptErr := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if bcryptErr != nil {
		log.Println("Unable to Encrypt Password:", bcryptErr)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Error While Generating Password, Try again later"},
		})
	}

	user.Password = string(bcryptPassword)
	dbErr = u.Store.SavePassword(*user)
	if dbErr != nil {
		log.Println("Error while updating password", dbErr)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Error while Updating Password"},
		})
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   &APISuccessResponse{Detail: "Password Reset Successfully"},
	})
}
