package neteasecloudmusic

func initRoute() {
	f := []requestFunc{
		Login,
		LoginCellphone,
		PlaylistDetail,
		SongDetail,
		SongUrl,
		ActivateInitProfile,
		Album,
		AlbumDetailDynamic,
		AlbumNewest,
		AlbumSub,
		AlbumSublist,
		ArtistAlbum,
		ArtistDesc,
		ArtistList,
		ArtistMv,
		ArtistSub,
		ArtistSublist,
		ArtistTopSong,
		Artists,
		Banner,
		Batch,
		CaptchaSent,
		CaptchaVerify,
		CellphoneExistenceCheck,
		CheckMusic,
		Comment,
		CommentAlbum,
		CommentDj,
		CommentEvent,
		CommentHot,
		CommentHotwallList,
		CommentLike,
		CommentMusic,
		Search,
	}
	route = funcNames(f)
	urls = make([]string, 0)
	for k, _ := range route {
		urls = append(urls, k+"\n")
	}
}
