package main

import (
	// "encoding/base64"
	// "fmt"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	// "io/fs"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	const maxMemory = 10 << 20 // 10 MB
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()

	mediaType := header.Header.Get("Content-Type")
	if mediaType == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Content-Type for thumbnail", nil)
		return
	}

	// data, err := io.ReadAll(file)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Error reading file", err)
	// 	return
	// }

	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find video", err)
		return
	}
	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Not authorized to update this video", nil)
		return
	}

	// videoThumbnails[videoID] = thumbnail{
	// 	data:      data,
	// 	mediaType: mediaType,
	// }

	// encodedDat := base64.StdEncoding.EncodeToString(data)
	// dataUrl := fmt.Sprintf("data:%v;base64,%v", mediaType, encodedDat)
	extension := strings.Split(mediaType, "/")[1]

	target := make([]byte, 32)
	rand.Read(target)
	thumbnailName := base64.RawURLEncoding.EncodeToString(target)
	fileName := fmt.Sprintf("/mnt/c/projects/cdn_go/assets/%v.%v", thumbnailName, extension)

	// create a new file
	filePath, err := os.Create(fileName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file", err)
		return
	}
	_, err = io.Copy(filePath, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to copy to file", err)
		return
	}

	thumbnailUrl := fmt.Sprintf("http://localhost:%v/assets/%v.%v", cfg.port, thumbnailName, extension)
	if mediaType == "" {
		respondWithError(w, http.StatusBadRequest, "Thumbnail Url missing", nil)
		return
	}

	video.ThumbnailURL = &thumbnailUrl

	err = cfg.db.UpdateVideo(video)
	if err != nil {
		delete(videoThumbnails, videoID)
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
