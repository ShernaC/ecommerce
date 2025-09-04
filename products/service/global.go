package service

import (
	"fmt"
	"products/config"
	"runtime/debug"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func GetService() *Service {
	s := Service{
		DB: config.GetDB(),
	}
	return &s
}

func GetTransaction() *Service {
	fmt.Println("begin...")
	s := Service{
		DB: config.GetDB().Begin(),
	}
	return &s
}

func (s *Service) Commit() error {
	if err := s.DB.Commit().Error; err != nil {
		return err
	}

	fmt.Println("commit...")

	return nil
}

func (s *Service) Rollback(err ...interface{}) error {
	s.DB.Rollback()
	fmt.Println("rollback...")
	fmt.Println("panic error: ", err)
	debug.PrintStack()

	if len(err) > 0 && err[0] != nil {
		return err[0].(error)
	}

	return fmt.Errorf("error")
}

func (s *Service) ErrorCheck(err ...interface{}) error {
	fmt.Println("panic error:", err)
	debug.PrintStack()

	if len(err) > 0 && err[0] != nil {
		return err[0].(error)
	}

	return nil
}
