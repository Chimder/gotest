package server

import (
	"context"
	"log"
	"net"

	"github.com/chimas/GoProject/internal/repository"
	"github.com/chimas/GoProject/proto/chapter"
	"github.com/chimas/GoProject/proto/manga"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func StartGRPCServer(db *pgxpool.Pool, rdb *redis.Client) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("err listen")
	}

	repo := repository.NewRepository(db)

	grpcServer := grpc.NewServer()

	manga.RegisterMangaServiceServer(grpcServer, NewMangaProto(repo))
	chapter.RegisterChapterServiceServer(grpcServer, NewChapterProto(repo))

	reflection.Register(grpcServer)

	log.Println("gRPC is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type MangaProto struct {
	manga.UnimplementedMangaServiceServer
	repo *repository.Repository
}

func NewMangaProto(repo *repository.Repository) *MangaProto {
	return &MangaProto{
		repo: repo,
	}
}

func (mp *MangaProto) GetManga(ctx context.Context, req *manga.MangaRequest) (*manga.MangaResponseWithChapters, error) {
	r, err := mp.repo.Manga.GetMangaByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	ch, err := mp.repo.Chapter.ListChaptersByMangaName(ctx, req.Name)

	chapters := make([]*chapter.ChaptersResponse, 0, len(ch))
	for _, c := range ch {
		chapters = append(chapters, &chapter.ChaptersResponse{
			Chapter:   c.Chapter,
			Img:       c.Img,
			Name:      c.Name,
			MangaName: c.AnimeName,
			CreatedAt: timestamppb.New(c.CreatedAt),
		})
	}

	return &manga.MangaResponseWithChapters{
		Manga: &manga.MangaResponse{
			Id:            r.ID,
			Name:          r.Name,
			Img:           r.Img,
			ImgHeader:     r.ImgHeader,
			Describe:      r.Describe,
			Genres:        r.Genres,
			Author:        r.Author.String,
			Country:       r.Country,
			Published:     int32(r.Published),
			AverageRating: r.AverageRating,
			RatingCount:   int32(r.RatingCount),
			Status:        r.Status,
			Popularity:    int32(r.Popularity),
		}, Chapters: chapters,
	}, nil
}

func (mp *MangaProto) GetAllMangas(ctx context.Context, req *emptypb.Empty) (*manga.MangaListResponse, error) {
	r, err := mp.repo.Manga.ListMangas(ctx)
	if err != nil {
		return nil, err
	}
	mangaList := make([]*manga.MangaResponse, 0, len(r))

	for _, m := range r {
		mangaList = append(mangaList, &manga.MangaResponse{
			Id:            m.ID,
			Name:          m.Name,
			Img:           m.Img,
			ImgHeader:     m.ImgHeader,
			Describe:      m.Describe,
			Genres:        m.Genres,
			Author:        m.Author.String,
			Country:       m.Country,
			Published:     int32(m.Published),
			AverageRating: m.AverageRating,
			RatingCount:   int32(m.RatingCount),
			Status:        m.Status,
			Popularity:    int32(m.Popularity),
		})
	}

	return &manga.MangaListResponse{
		MangaLists: mangaList,
	}, nil
}
func (mp *MangaProto) GetPopularMangas(ctx context.Context, req *emptypb.Empty) (*manga.MangaListResponse, error) {
	r, err := mp.repo.Manga.ListPopularMangas(ctx)
	if err != nil {
		return nil, err
	}

	mangalist := make([]*manga.MangaResponse, 0, len(r))
	for _, m := range r {
		mangalist = append(mangalist, &manga.MangaResponse{
			Id:            m.ID,
			Name:          m.Name,
			Img:           m.Img,
			ImgHeader:     m.ImgHeader,
			Describe:      m.Describe,
			Genres:        m.Genres,
			Author:        m.Author.String,
			Country:       m.Country,
			Published:     int32(m.Published),
			AverageRating: m.AverageRating,
			RatingCount:   int32(m.RatingCount),
			Status:        m.Status,
			Popularity:    int32(m.Popularity),
		})
	}
	return &manga.MangaListResponse{MangaLists: mangalist}, nil
}

func (mp *MangaProto) GetFilteredMangas(ctx context.Context, req *manga.MangaFilterRequest) (*manga.MangaListResponse, error) {
	return nil, nil
}

type ChapterProto struct {
	chapter.UnimplementedChapterServiceServer
	repo *repository.Repository
}

func NewChapterProto(r *repository.Repository) *ChapterProto {
	return &ChapterProto{
		repo: r,
	}
}

func (ch *ChapterProto) GetChapters(ctx context.Context, req *chapter.ChaptersRequest) (*chapter.ChaptersResponse, error) {
	c, err := ch.repo.Chapter.GetChapterByMangaNameAndNumber(ctx, req.Name, int(req.Chapter))
	if err != nil {
		return nil, err
	}
	return &chapter.ChaptersResponse{
		Chapter:   c.Chapter,
		Img:       c.Img,
		Name:      c.Name,
		MangaName: c.AnimeName,
		CreatedAt: timestamppb.New(c.CreatedAt),
	}, nil
}
