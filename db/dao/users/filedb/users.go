package usersdao

import (
	// "fmt"
	"context"

	"github.com/ANMalko/grpc-server.git/db/model"
	"github.com/ANMalko/grpc-server.git/db/filedb"
	"github.com/ANMalko/grpc-server.git/db/error"
)

type DAO struct {
	fileDB *filedb.FileDB
}

func (d *DAO) DB() *filedb.FileDB{
	return d.fileDB
}

func NewDAO(ctx context.Context, fileName string) *DAO{
	db := filedb.NewFileDB(ctx, fileName)
	dao := DAO{fileDB: db}
	return &dao

}

func (f *DAO) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {

	createdUser := model.User{
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
		PhoneNumber:user.PhoneNumber,
	}

	if ok := f.fileDB.AddUser(&createdUser); !ok {
		err := dberror.New(dberror.EALREADYEXISTS)
		return nil, err

	}
	return &createdUser, nil
}

func (f *DAO) UpdateUser(ctx context.Context, user *model.User) error {
	if ok := f.fileDB.UpdateUser(user); !ok {
		err := dberror.New(dberror.ENOTFOUND)
		return err
	}
	return nil
}

func (f *DAO) GetUser(ctx context.Context, userID uint32) *model.User {
	user := f.fileDB.GetUser(userID)

	return user
}

func (f *DAO) DeleteUser(ctx context.Context, userID uint32) {
	f.fileDB.DeleteUser(userID)
}
