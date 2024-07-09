package postgres

import (
	pb "auth-service/generated/auth_service"
	"database/sql"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (u *UserRepo) Register(user *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var userId string
	err := u.DB.QueryRow(`
		INSERT INTO users (
			username,
			password,
			email
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			id
	`, user.Username, user.Password, user.Email).Scan(&userId)

	if err != nil {
		return &pb.RegisterResponse{
			Message: "Failed to create user",
			UserId:  "",
		}, nil
	}

	return &pb.RegisterResponse{
		Message: "User created successfully",
		UserId:  userId,
	}, nil
}

func (u *UserRepo) Login(login *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user pb.LoginResponse

	err := u.DB.QueryRow(`
		SELECT
			id,
			username,
			email
		FROM
			users
		WHERE
			username = $1 and password = $2 and deleted_at = 0
	`, login.Username, login.Password).Scan(&user.UserId, &user.Username, &user.Email)

	return &user, err
}

func (u *UserRepo) LogoutUser(id string) (*pb.LogoutResponse, error) {
	_, err := u.DB.Exec(`
		UPDATE 
			users 
		SET
			deleted_at = DATE_PART('epoch', CURRENT_TIMESTAMP)::INT
		WHERE
			deleted_at = 0 and id = $1
	`, id)

	if err != nil {
		return &pb.LogoutResponse{
			Message: "Faild to deleted user",
		}, err
	}

	return &pb.LogoutResponse{
		Message: "user deleted successully",
	}, nil
}

func (u *UserRepo) CreateProfile(profile *pb.UpdateUserProfileRequest) error {
	_, err := u.DB.Exec(`
		INSERT INTO user_profiles (
			user_id,
			username,
			fullname,
			date_of_birth,
			phone_number,
			address
		) 
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
	`, profile.UserId, profile.Username, profile.FullName, profile.DateOfBirth, profile.PhoneNumber, profile.Address)

	return err
}

func (u *UserRepo) GetUserProfile(username string) (*pb.GetUserProfileResponse, error) {
	var userProfile pb.GetUserProfileResponse
	err := u.DB.QueryRow(`
		SELECT
			fullname,
			username,
			date_of_birth,
			phone_number,
			address
		FROM
			user_profiles
		WHERE
			username = $1
	`, username).Scan(&userProfile.Fullname, &userProfile.Username, &userProfile.DateOfBirth, &userProfile.PhoneNumber, &userProfile.PhoneNumber)

	return &userProfile, err
}

func (u *UserRepo) UpdateUserProfile(profile *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	_, err := u.DB.Exec(`
		UPDATE 
			user_profiles
		SET
			fullname = $1,
			username = $2,
			date_of_birth = $3,
			phone_number = $4,
			address = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			user_id = $6
	`, profile.FullName, profile.Username, profile.DateOfBirth, profile.PhoneNumber, profile.Address, profile.UserId)

	if err != nil {
		return &pb.UpdateUserProfileResponse{
			Message: "Faild to updated user",
		}, err
	}

	return &pb.UpdateUserProfileResponse{
		Message: "User updated successfully",
	}, nil
}
