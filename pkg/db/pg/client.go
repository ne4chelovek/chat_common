package pg

import (
	"github.com/ne4chelovek/chat_common/pkg/db"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// pgClient - структура, реализующая интерфейс db.Client.
// Содержит masterDBC, который представляет собой соединение с базой данных.
type pgClient struct {
	masterDBC db.DB // Интерфейс для работы с базой данных
}

// New - функция-конструктор для создания нового клиента базы данных.
// Принимает контекст (ctx) и строку подключения (dsn).
// Возвращает интерфейс db.Client и ошибку, если что-то пошло не так.
func New(ctx context.Context, dsn string) (db.Client, error) {
	// Создаем пул соединений с базой данных PostgreSQL
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		// Если произошла ошибка, возвращаем её с описанием
		return nil, fmt.Errorf("failed to connect to db: %v", err)
	}

	// Возвращаем новый экземпляр pgClient, где masterDBC инициализирован
	// с использованием структуры pg, которая реализует интерфейс db.DB.
	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

// DB - метод, возвращающий интерфейс db.DB для работы с базой данных.
// Это позволяет клиенту получить доступ к методам базы данных.
func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

// Close - метод для закрытия соединения с базой данных.
// Если masterDBC не равен nil, вызывается метод Close() для закрытия соединения.
// Возвращает ошибку, если что-то пошло не так.
func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		// Закрываем соединение с базой данных
		c.masterDBC.Close()
	}

	return nil
}
