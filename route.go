package hammer

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
	}
	route = funcNames(f)
}
