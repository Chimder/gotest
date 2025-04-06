package service

import (
	"context"
	"database/sql"
	"fmt"

	"slices"

	"github.com/chimas/GoProject/internal/models"
	"github.com/chimas/GoProject/internal/repository"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	repo *repository.Repository
	rdb  *redis.Client
}

func NewUserService(repo *repository.Repository, rdb *redis.Client) *UserService {
	return &UserService{repo: repo, rdb: rdb}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.UserRepo, error) {
	return s.repo.User.GetUserByEmail(ctx, email)
}

func (s *UserService) InsertUser(ctx context.Context, arg *models.UserRepo) (models.UserRepo, error) {
	user, err := s.repo.User.GetUserByEmail(ctx, arg.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			user, err = s.repo.User.InsertUser(ctx, repository.InsertUserParams{
				Email: arg.Email,
				Name:  arg.Name,
				Image: arg.Image,
			})
			if err != nil {
				return models.UserRepo{}, fmt.Errorf("no rows %v", err)
			}
		} else {
			return models.UserRepo{}, err
		}
	}
	return user, nil
}

func (s *UserService) ToggleFavorite(ctx context.Context, email string, name string) error {

	user, err := s.repo.User.GetUserByEmail(ctx, email)
	if err != nil {
		// utils.WriteError(w, 500, op+"GUBE", err)
		return err
	}

	isAnimeInFavorites := slices.Contains(user.Favorite, name)

	if !isAnimeInFavorites {
		err = s.repo.Manga.UpdateMangaPopularity(ctx, name)
		if err != nil {
			return err
		}

		user.Favorite = append(user.Favorite, name)
		err = s.repo.User.UpdateUserFavorites(ctx, user.Favorite, email)
		if err != nil {
			return err
		}

		return nil
		// if err := json.NewEncoder(w).Encode(SuccessResponse{Success: "Manga added"}); err != nil {
		// 	utils.WriteError(w, 500, op+"ENC", err)
		// 	return
		// }
	} else {
		newFavorites := []string{}
		for _, favorite := range user.Favorite {
			if favorite != name {
				newFavorites = append(newFavorites, favorite)
			}
		}
		user.Favorite = newFavorites
		err = s.repo.User.UpdateUserFavorites(ctx, newFavorites, email)
		if err != nil {
			return err
		}
		return nil

		// if err := utils.WriteJSON(w, 200, SuccessResponse{Success: "Manga deleted"}); err != nil {
		// 	utils.WriteError(w, 500, op+"WJ", err)
		// 	return
		// }
	}
}
func (s *UserService) GetUserFavorites(ctx context.Context, email string) ([]models.MangaRepo, error) {
	favorites, err := s.repo.User.GetUserFavoritesByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if len(favorites) == 0 {
		return nil, fmt.Errorf("len favorite zero %v", err)
	}

	return s.repo.Manga.GetMangaByNames(ctx, favorites)
}
func (s *UserService) IsUserFavorite(ctx context.Context, email string, name string) (bool, error) {
	user, err := s.repo.User.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	isAnimeInFavorites := slices.Contains(user.Favorite, name)
	return isAnimeInFavorites, nil
}

func (s *UserService) UpdateUserFavorites(ctx context.Context, favorites []string, email string) error {
	return s.repo.User.UpdateUserFavorites(ctx, favorites, email)
}

func (s *UserService) DeleteUser(ctx context.Context, email string) error {
	return s.repo.User.DeleteUserByEmail(ctx, email)
}
