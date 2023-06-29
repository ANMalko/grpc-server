package filedb

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/ANMalko/grpc-server.git/db/model"
)

type FileDB struct {
	fileName string
	dB map[uint32]*model.User
}

func NewFileDB(ctx context.Context, fileName string) *FileDB{
	dB := make(map[uint32]*model.User, 10)
	fileDB := FileDB{fileName, dB}
	fileDB.loadDB()
	return &fileDB
}

func (f *FileDB) AddUser(user *model.User) bool {
	if _, ok := f.dB[user.Id]; ok {
		return false
	}
	f.dB[user.Id] = user
	return true
}

func (f *FileDB) GetUser(userId uint32) *model.User {
	var user *model.User
	user, ok := f.dB[userId]
	if !ok {
		return nil
	}
	return user
}

func (f *FileDB) DeleteUser(userId uint32) {
	delete(f.dB, userId)
}

func (f *FileDB) UpdateUser(user *model.User) bool {
	_, ok := f.dB[user.Id]
	if !ok {
		return false
	}
	f.dB[user.Id] = user
	return true
}

func (f *FileDB) loadDB() {
	file, err := os.Open(f.fileName)
	if err != nil {
		return
    }
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var user *model.User
	for scanner.Scan() {
		json.Unmarshal([]byte(scanner.Text()), &user)
		f.dB[user.Id] = user
	}
}

func (f *FileDB) DumpDB() {
	log.Info().Msg("Dump userBD to file")

	if len(f.dB) == 0 {
		log.Info().Msg("DB is empty")
		return
	}

	file, err := os.Create(f.fileName)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to save userDB to %s", f.fileName)
		return
    }
	defer file.Close()

	for userId := range f.dB {
		userDumped, err := json.Marshal(*f.dB[userId])
        if err != nil {
            return
        }
		line := fmt.Sprintf("%s\n", string(userDumped))
	    file.WriteString(line)
	}
}
