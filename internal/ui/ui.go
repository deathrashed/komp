package ui

import "errors"

func Interactive() bool { return false }

func PickFiles(dir string) ([]string, error) { return nil, errors.New("stub") }

func PickFormat() (string, error) { return "", errors.New("stub") }

func ConfirmDelete() bool { return false }
