package repository

import (
	"database/sql"
	"time"

	"github.com/ombima56/insights-edge/backend/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByWallet(walletAddr string) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
}

type InsightRepository interface {
	CreateInsight(insight *models.Insight) error
	GetInsight(id int64) (*models.Insight, error)
	GetInsights() ([]*models.Insight, error)
	CreatePurchase(purchase *models.Purchase) error
	GetPurchasesByWallet(walletAddr string) ([]*models.Purchase, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (email, password, wallet_addr, first_name, last_name, account_type, company_name, industry, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, user.Email, user.Password, user.WalletAddr, user.FirstName, user.LastName, user.AccountType, user.CompanyName, user.Industry, time.Now().Format(time.RFC3339))
	return err
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	row := r.db.QueryRow(`
		SELECT id, email, password, wallet_addr, first_name, last_name, account_type, company_name, industry, created_at
		FROM users
		WHERE email = ?
	`, email)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.WalletAddr, &user.FirstName, &user.LastName, &user.AccountType, &user.CompanyName, &user.Industry, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) GetUserByWallet(walletAddr string) (*models.User, error) {
	var user models.User
	row := r.db.QueryRow(`
		SELECT id, email, password, wallet_addr, first_name, last_name, account_type, company_name, industry, created_at
		FROM users
		WHERE wallet_addr = ?
	`, walletAddr)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.WalletAddr, &user.FirstName, &user.LastName, &user.AccountType, &user.CompanyName, &user.Industry, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) GetUserByID(id int64) (*models.User, error) {
	var user models.User
	row := r.db.QueryRow(`
		SELECT id, email, password, wallet_addr, first_name, last_name, account_type, company_name, industry, created_at
		FROM users
		WHERE id = ?
	`, id)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.WalletAddr, &user.FirstName, &user.LastName, &user.AccountType, &user.CompanyName, &user.Industry, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

type insightRepository struct {
	db *sql.DB
}

func NewInsightRepository(db *sql.DB) InsightRepository {
	return &insightRepository{db: db}
}

func (r *insightRepository) CreateInsight(insight *models.Insight) error {
	_, err := r.db.Exec(`
		INSERT INTO insights (provider, industry, title, description, price, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, insight.Provider, insight.Industry, insight.Title, insight.Description, insight.Price, time.Now().Format(time.RFC3339))
	return err
}

func (r *insightRepository) GetInsight(id int64) (*models.Insight, error) {
	var insight models.Insight
	row := r.db.QueryRow(`
		SELECT id, provider, industry, title, description, price, created_at
		FROM insights
		WHERE id = ?
	`, id)

	err := row.Scan(&insight.ID, &insight.Provider, &insight.Industry, &insight.Title, &insight.Description, &insight.Price, &insight.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &insight, err
}

func (r *insightRepository) GetInsights() ([]*models.Insight, error) {
	rows, err := r.db.Query(`
		SELECT id, provider, industry, title, description, price, created_at
		FROM insights
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var insights []*models.Insight
	for rows.Next() {
		var insight models.Insight
		err := rows.Scan(&insight.ID, &insight.Provider, &insight.Industry, &insight.Title, &insight.Description, &insight.Price, &insight.CreatedAt)
		if err != nil {
			return nil, err
		}
		insights = append(insights, &insight)
	}
	return insights, nil
}

func (r *insightRepository) CreatePurchase(purchase *models.Purchase) error {
	_, err := r.db.Exec(`
		INSERT INTO purchases (insight_id, buyer, created_at)
		VALUES (?, ?, ?)
	`, purchase.InsightID, purchase.Buyer, time.Now().Format(time.RFC3339))
	return err
}

func (r *insightRepository) GetPurchasesByWallet(walletAddr string) ([]*models.Purchase, error) {
	rows, err := r.db.Query(`
		SELECT id, insight_id, buyer, created_at
		FROM purchases
		WHERE buyer = ?
		ORDER BY created_at DESC
	`, walletAddr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*models.Purchase
	for rows.Next() {
		var purchase models.Purchase
		err := rows.Scan(&purchase.ID, &purchase.InsightID, &purchase.Buyer, &purchase.CreatedAt)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, &purchase)
	}
	return purchases, nil
}
