package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
)

var DB *gorm.DB

func Connect() error {
	// Local
	// err := godotenv.Load()
	// if err != nil {
	// 	return fmt.Errorf("erro ao carregar o arquivo .env: %w", err)
	// }

	// Cloud render
	err := godotenv.Load()
	if err != nil {
		// Em vez de retornar erro e derrubar o app, apenas avisa no terminal
		fmt.Println("Aviso: Arquivo .env não encontrado. Lendo variáveis do sistema (Produção).")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=America/Sao_Paulo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("falha ao conectar no banco de dados: %w", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Character{},
		&models.Event{},
		&models.Guide{},
		&models.LastEvent{},
	)
	if err != nil {
		return fmt.Errorf("falha ao executar o AutoMigrate: %w", err)
	}

	// if err := createConstraints(); err != nil {
	// 	return fmt.Errorf("falha ao criar constraints: %w", err)
	// }

	fmt.Println("Conexão com o banco de dados e migrações concluídas com sucesso!")
	return nil
}

func createConstraints() error {
	constraints := []string{
		`DO $$ 
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_guides_user') THEN
            ALTER TABLE guides ADD CONSTRAINT fk_guides_user
                FOREIGN KEY (user_id) REFERENCES users(user_id)
                ON UPDATE CASCADE ON DELETE CASCADE;
        END IF;
    END $$;`,

		`DO $$ 
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_guides_character') THEN
            ALTER TABLE guides ADD CONSTRAINT fk_guides_character
                FOREIGN KEY (character_id) REFERENCES characters(character_id)
                ON UPDATE CASCADE ON DELETE CASCADE;
        END IF;
    END $$;`,

		`DO $$ 
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_events_user') THEN
            ALTER TABLE events ADD CONSTRAINT fk_events_user
                FOREIGN KEY (user_id) REFERENCES users(user_id)
                ON UPDATE CASCADE ON DELETE CASCADE;
        END IF;
    END $$;`,

		`DO $$ 
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_last_events_event') THEN
            ALTER TABLE last_events ADD CONSTRAINT fk_last_events_event
                FOREIGN KEY (event_id) REFERENCES events(event_id)
                ON UPDATE CASCADE ON DELETE CASCADE;
        END IF;
    END $$;`,
	}

	for _, f := range fks {
		if !m.HasConstraint(f.model, f.field) {
			if err := m.CreateConstraint(f.model, f.field); err != nil {
				return fmt.Errorf("erro criando constraint %s.%s: %w", fmt.Sprintf("%T", f.model), f.field, err)
			}
		}
	}
	return nil
}
