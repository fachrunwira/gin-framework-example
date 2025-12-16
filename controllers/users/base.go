package users

import "log/slog"

type userControllers struct {
	appLog *slog.Logger
}

func UserControllers(logger *slog.Logger) *userControllers {
	return &userControllers{
		appLog: logger,
	}
}
