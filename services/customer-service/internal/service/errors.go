package service

import "errors"

var (
	ErrNotFound         = errors.New("запись не найдена")
	ErrInvalidID        = errors.New("некорректный идентификатор")
	ErrEmptyComposition = errors.New("состав не может быть пустым")
	ErrNotEnoughStock   = errors.New("недостаточно товара на складе")
	ErrInvalidCount     = errors.New("количество должно быть больше 0")
	ErrExternalService  = errors.New("ошибка внешнего сервиса")
	ErrInvalidDate      = errors.New("некорректная дата")
)
