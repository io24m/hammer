package hammer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

const (
	post                       string = "POST"
	get                        string = "GET"
	urlLogin                   string = "https://music.163.com/weapi/login"
	urlLoginCellphone          string = "https://music.163.com/weapi/login/cellphone"
	urlPlaylistDetail          string = "https://music.163.com/weapi/v3/playlist/detail"
	urlSongDetail              string = "https://music.163.com/weapi/v3/song/detail"
	urlSongUrl                 string = "https://music.163.com/api/song/enhance/player/url"
	urlActivateInitProfile     string = "http://music.163.com/eapi/activate/initProfile"
	urlAlbum                   string = "https://music.163.com/weapi/v1/album/%s"
	urlAlbumDetailDynamic      string = "https://music.163.com/api/album/detail/dynamic"
	urlAlbumNewest             string = "https://music.163.com/api/discovery/newAlbum"
	urlAlbumSub                string = "https://music.163.com/api/album/%s"
	urlAlbumSublist            string = "https://music.163.com/weapi/album/sublist"
	urlArtistAlbum             string = "https://music.163.com/weapi/artist/albums/%s"
	urlArtistDesc              string = "https://music.163.com/weapi/artist/introduction"
	urlArtistList              string = "https://music.163.com/api/v1/artist/list"
	urlArtistMv                string = "https://music.163.com/weapi/artist/mvs"
	urlArtistSub               string = "https://music.163.com/weapi/artist/%s"
	urlArtistSublist           string = "https://music.163.com/weapi/artist/sublist"
	urlArtistTopSong           string = "https://music.163.com/api/artist/top/song"
	urlArtists                 string = "https://music.163.com/weapi/v1/artist/%s"
	urlBanner                  string = "https://music.163.com/api/v2/banner/get"
	urlBatch                   string = "http://music.163.com/eapi/batch"
	urlCaptchaSent             string = "https://music.163.com/weapi/sms/captcha/sent"
	urlCaptchaVerify           string = "https://music.163.com/weapi/sms/captcha/verify"
	urlCellphoneExistenceCheck string = "http://music.163.com/eapi/cellphone/existence/check"
	urlCheckMusic              string = "https://music.163.com/weapi/song/enhance/player/url"
	urlComment                 string = "https://music.163.com/weapi/resource/comments/%s"
	urlCommentAlbum            string = "https://music.163.com/weapi/v1/resource/comments/R_AL_3_%s"
	urlCommentDj               string = "https://music.163.com/weapi/v1/resource/comments/A_DJ_1_%s"
	urlCommentEvent            string = "https://music.163.com/weapi/v1/resource/comments/%s"
	urlCommentHot              string = "https://music.163.com/weapi/v1/resource/hotcomments/%s%s"
	urlCommentHotwallList      string = "https://music.163.com/api/comment/hotwall/list/get"
	urlCommentLike             string = "https://music.163.com/weapi/v1/comment/%s"
	urlCommentMusic            string = "https://music.163.com/api/v1/resource/comments/R_SO_4_%s"
)

//Login 邮箱登录
func Login(query *Query) (string, error) {
	var data = make(map[string]interface{})
	data["username"] = query.GetParam("email")
	data["rememberLogin"] = "true"
	if md5Password := query.GetParam("md5_password"); strings.TrimSpace(md5Password) != "" {
		data["password"] = md5Password
	} else {
		pw := query.GetParam("password")
		sum := md5.Sum([]byte(pw))
		data["password"] = hex.EncodeToString(sum[:])
	}
	query.AddCookie("os", "pc")
	var opt = &options{
		crypto:  weapi,
		ua:      pc,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	cmResult, err := requestCloudMusicApi(post, urlLogin, data, opt)
	if err != nil {
		return "", err
	}
	defer cmResult.Body.Close()
	body, err := ioutil.ReadAll(cmResult.Body)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{}, 0)
	json.Unmarshal(body, &m)
	float := strconv.FormatFloat(m["code"].(float64), 'f', -1, 64)
	if float == "502" {
		return "账号或密码错误", nil
	}
	m["cookie"] = cmResult.Header.Get("Set-Cookie")
	marshal, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	//code502	var msg = `{'msg':'账号或密码错误','code':'502','message':'账号或密码错误'}`
	return string(marshal), nil
}

//LoginCellphone 电话登录
func LoginCellphone(query *Query) (string, error) {
	data := make(map[string]interface{}, 0)
	data["phone"] = query.GetParam("phone")
	if cc := query.GetParam("countrycode"); strings.TrimSpace(cc) != "" {
		data["countrycode"] = query.GetParam("countrycode")
	}
	data["rememberLogin"] = "true"
	if md5Password := query.GetParam("md5_password"); strings.TrimSpace(md5Password) != "" {
		data["password"] = md5Password
	} else {
		pw := query.GetParam("password")
		sum := md5.Sum([]byte(pw))
		data["password"] = hex.EncodeToString(sum[:])
	}
	query.AddCookie("os", "pc")
	opt := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   nil,
		ua:      pc,
	}
	cmResult, err := requestCloudMusicApi(post, urlLoginCellphone, data, opt)
	if err != nil {
		return "", err
	}
	defer cmResult.Body.Close()
	body, err := ioutil.ReadAll(cmResult.Body)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{}, 0)
	json.Unmarshal(body, &m)
	m["cookie"] = cmResult.Header.Get("Set-Cookie")
	marshal, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

//PlaylistDetail 歌单详情
func PlaylistDetail(query *Query) (string, error) {
	data := make(map[string]interface{}, 0)
	data["id"] = query.GetParam("id")
	data["n"] = 100000
	data["s"] = 8
	opt := &options{
		crypto:  linuxapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlPlaylistDetail, data, opt)
}

//SongUrl 歌曲链接
func SongUrl(query *Query) (string, error) {
	if MUSIC_U := query.GetCookie("MUSIC_U"); strings.TrimSpace(MUSIC_U) == "" {
		query.AddCookie("_ntes_nuid", hex.EncodeToString(key(16)))
	}
	query.AddCookie("os", "pc")
	data := make(map[string]interface{}, 0)
	data["ids"] = "[" + query.GetParam("id") + "]"
	data["br"] = query.GetParamOrDefault("br", 999000)
	opt := &options{
		crypto:  linuxapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlSongUrl, data, opt)
}

//SongDetail 歌曲详情
func SongDetail(query *Query) (string, error) {
	ids := query.GetParam("ids")
	reg, _ := regexp.Compile(`\s*,\s*`)
	idList := reg.Split(ids, -1)
	c := make([]string, 0)
	for _, v := range idList {
		c = append(c, fmt.Sprintf(`{"id":%s}`, v))
	}
	data := make(map[string]interface{})
	data["c"] = "[" + strings.Join(c, ",") + "]"
	data["ids"] = "[" + strings.Join(idList, ",") + "]"
	opt := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlSongDetail, data, opt)
}

//ActivateInitProfile 初始化名字
func ActivateInitProfile(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["nickname"] = query.GetParam("nickname")
	opt := &options{
		crypto:  eapi,
		cookies: query.Cookies,
		url:     "/api/activate/initProfile",
	}
	return responseDefault(post, urlActivateInitProfile, data, opt)
}

//Album 专辑内容
func Album(query *Query) (string, error) {
	id := query.GetParam("id")
	opt := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, fmt.Sprintf(urlAlbum, id), nil, opt)
}

//AlbumDetailDynamic 专辑动态信息
func AlbumDetailDynamic(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["id"] = query.GetParam("id")
	opt := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlAlbumDetailDynamic, data, opt)
}

//AlbumNewest 最新专辑
func AlbumNewest(query *Query) (string, error) {
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlAlbumNewest, nil, opt)
}

//AlbumSub 收藏/取消收藏专辑
func AlbumSub(query *Query) (string, error) {
	data := make(map[string]interface{})
	t := query.GetParam("t")
	data["id"] = query.GetParam("id")
	if t == "1" {
		t = "sub"
	} else {
		t = "unsub"
	}
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlAlbumSub, t), data, opt)
}

//AlbumSublist 已收藏专辑列表
func AlbumSublist(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["total"] = true
	data["limit"] = query.GetParamOrDefault("limit", 25)
	data["offset"] = query.GetParamOrDefault("offset", 0)
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlAlbumSublist, data, opt)
}

//ArtistAlbum 歌手专辑列表
func ArtistAlbum(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["total"] = true
	data["limit"] = query.GetParamOrDefault("limit", 30)
	data["offset"] = query.GetParamOrDefault("offset", 0)
	id := query.GetParam("id")
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlArtistAlbum, id), data, opt)
}

//ArtistDesc 歌手介绍
func ArtistDesc(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["id"] = query.GetParam("id")
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlArtistDesc, data, opt)
}

//ArtistList 歌手分类
//type 取值[1:男歌手、2:女歌手、3:乐队]
//area 取值[-1:全部、7华语、96欧美、8:日本、16韩国、0:其他]
//initial 取值 a-z/A-Z
func ArtistList(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["initial"] = query.GetParam("initial")
	data["offset"] = query.GetParamOrDefault("offset", 0)
	data["limit"] = query.GetParamOrDefault("limit", 30)
	data["total"] = true
	data["type"] = query.GetParamOrDefault("type", "1")
	data["area"] = query.GetParam("area")
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlArtistList, data, opt)
}

//ArtistMv 歌手相关MV
func ArtistMv(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["artistId"] = query.GetParam("id")
	data["limit"] = query.GetParam("limit")
	data["offset"] = query.GetParam("offset")
	data["total"] = true
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlArtistMv, data, opt)
}

//ArtistSub 收藏与取消收藏歌手
func ArtistSub(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["artistId"] = query.GetParam("id")
	data["artistIds"] = "[" + query.GetParam("id") + "]"
	t := "sub"
	if query.GetParam("t") != "1" {
		t = "unsub"
	}
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlArtistSub, t), data, opt)
}

//ArtistSublist 关注歌手列表
func ArtistSublist(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["limit"] = query.GetParamOrDefault("limit", 25)
	data["offset"] = query.GetParamOrDefault("offset", 0)
	data["total"] = true
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlArtistSublist, data, opt)
}

//ArtistTopSong 歌手热门 50 首歌曲
func ArtistTopSong(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["id"] = query.GetParam("id")
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlArtistTopSong, data, opt)
}

//Artists 歌手单曲
func Artists(query *Query) (string, error) {
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlArtists, query.GetParam("id")), nil, opt)
}

//Banner 首页轮播图
func Banner(query *Query) (string, error) {
	platform := map[interface{}]string{
		"0": "pc",
		"1": "android",
		"2": "iphone",
		"3": "ipad",
	}
	data := make(map[string]interface{})
	t, ok := platform[query.GetParamOrDefault("type", "0")]
	if !ok {
		t = "pc"
	}
	data["clientType"] = t
	opt := &options{crypto: linuxapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlBanner, data, opt)
}

//Batch 批量请求接口
func Batch(query *Query) (string, error) {
	data := make(map[string]interface{})
	reg, _ := regexp.Compile(`^/api/`)
	data["e_r"] = true
	for k, v := range query.Param {
		if !reg.MatchString(k) {
			continue
		}
		data[k] = v
	}
	opt := &options{crypto: eapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlBatch, data, opt)
}

//CaptchaSent 发送验证码
func CaptchaSent(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["ctcode"] = query.GetParamOrDefault("ctcode", "86")
	data["cellphone"] = query.GetParam("phone")
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlCaptchaSent, data, opt)
}

//CaptchaVerify 校验验证码
func CaptchaVerify(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["ctcode"] = query.GetParamOrDefault("ctcode", "86")
	data["cellphone"] = query.GetParam("phone")
	data["captcha"] = query.GetParam("captcha")
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlCaptchaVerify, data, opt)
}

//CellphoneExistenceCheck 检测手机号码是否已注册
func CellphoneExistenceCheck(query *Query) (string, error) {
	data := map[string]interface{}{
		"cellphone":   query.GetParam("phone"),
		"countrycode": query.GetParam("countrycode"),
	}
	opt := &options{crypto: eapi, cookies: query.Cookies, proxy: query.Proxy, url: "/api/cellphone/existence/check"}
	return responseDefault(post, urlCellphoneExistenceCheck, data, opt)
}

//CheckMusic 歌曲可用性
func CheckMusic(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["ids"] = "[" + query.GetParam("id") + "]"
	data["br"] = query.GetParamOrDefault("br", 999000)
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	r, err := responseDefault(post, urlCheckMusic, data, opt)
	if err != nil {
		return "", err
	}
	j, err := ReadJson(r)
	if err != nil {
		return "", err
	}
	c, err := j.Get("code").Int()
	if err != nil {
		return "", err
	}
	if c == 200 {
		c, err = j.Get("data[0].code").Int()
		if err != nil {
			return "", err
		}
		if c == 200 {
			return `{success: true, message: 'ok'}`, nil
		}
	}
	return `{success: false, message: '暂无版权'}`, nil
}

//Comment 发送与删除评论
func Comment(query *Query) (string, error) {
	query.AddCookie("os", "pc")
	data := make(map[string]interface{})
	m := map[string]string{
		"0": "delete",
		"1": "add",
		"2": "reply",
	}
	t := m[query.GetParam("t")]
	m = map[string]string{
		"0": "R_SO_4_",  //歌曲
		"1": "R_MV_5_",  //MV
		"2": "A_PL_0_",  //歌单
		"3": "R_AL_3_",  //专辑
		"4": "A_DJ_1_",  //电台,
		"5": "R_VI_62_", //视频
		"6": "A_EV_2_",  //动态
	}
	tp := m[query.GetParam("type")]
	data["threadId"] = tp + query.GetParam("id")
	if tp == "A_EV_2_" {
		data["threadId"] = query.GetParam("threadId")
	}
	switch t {
	case "add":
		data["content"] = query.GetParam("content")
	case "delete":
		data["commentId"] = query.GetParam("commentId")
	case "reply":
		data["commentId"] = query.GetParam("commentId")
		data["content"] = query.GetParam("content")
	}
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlComment, t), data, opt)
}

//CommentAlbum 专辑评论
func CommentAlbum(query *Query) (string, error) {
	query.AddCookie("os", "pc")
	data := make(map[string]interface{})
	id := query.GetParam("id")
	data["rid"] = id
	data["limit"] = query.GetParamOrDefault("limit", 20)
	data["offset"] = query.GetParamOrDefault("offset", 0)
	data["beforeTime"] = query.GetParamOrDefault("before", 0)
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlCommentAlbum, id), data, opt)
}

//CommentDj 电台评论
func CommentDj(query *Query) (string, error) {
	query.AddCookie("os", "pc")
	data := make(map[string]interface{})
	id := query.GetParam("id")
	data["rid"] = id
	data["limit"] = query.GetParamOrDefault("limit", 20)
	data["offset"] = query.GetParamOrDefault("offset", 0)
	data["beforeTime"] = query.GetParamOrDefault("before", 0)
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlCommentDj, id), data, opt)
}

//CommentEvent 获取动态评论
func CommentEvent(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["limit"] = query.GetParamOrDefault("limit", 20)
	data["offset"] = query.GetParamOrDefault("offset", 0)
	data["beforeTime"] = query.GetParamOrDefault("before", 0)
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlCommentEvent, query.GetParam("threadId")), data, opt)
}

//CommentHot 热门评论
func CommentHot(query *Query) (string, error) {
	query.AddCookie("os", "pc")
	m := map[string]string{
		"0": "R_SO_4_",  //  歌曲
		"1": "R_MV_5_",  //  MV
		"2": "A_PL_0_",  //  歌单
		"3": "R_AL_3_",  //  专辑
		"4": "A_DJ_1_",  //  电台,
		"5": "R_VI_62_", //  视频
	}
	t := m[query.GetParam("type")]
	data := make(map[string]interface{})
	id := query.GetParam("id")
	data["rid"] = id
	data["limit"] = query.GetParamOrDefault("limit", 20)
	data["offset"] = query.GetParamOrDefault("offset", 0)
	data["beforeTime"] = query.GetParamOrDefault("before", 0)
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlCommentHot, t, id), data, opt)
}

//CommentHotwallList 云村热评
func CommentHotwallList(query *Query) (string, error) {
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlCommentHotwallList, nil, opt)
}

//CommentLike 点赞与取消点赞评论
func CommentLike(query *Query) (string, error) {
	query.AddCookie("os", "pc")
	t := "unlike"
	if query.GetParam("t") == "1" {
		t = "like"
	}
	m := map[string]string{
		"0": "R_SO_4_",  //  歌曲
		"1": "R_MV_5_",  //  MV
		"2": "A_PL_0_",  //  歌单
		"3": "R_AL_3_",  //  专辑
		"4": "A_DJ_1_",  //  电台,
		"5": "R_VI_62_", //  视频
		"6": "A_EV_2_",  //  动态
	}
	tp := m[query.GetParam("type")]
	data := make(map[string]interface{})
	data["threadId"] = tp + query.GetParam("id")
	if tp == "A_EV_2_" {
		data["threadId"] = query.GetParam("threadId")
	}
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlCommentLike, t), data, opt)
}

//CommentMusic 歌曲评论
func CommentMusic(query *Query) (string, error) {
	query.AddCookie("os", "pc")
	data := make(map[string]interface{})
	id := query.GetParam("id")
	data["rid"] = id
	data["limit"] = query.GetParamOrDefault("limit", 20)
	data["offset"] = query.GetParamOrDefault("offset", 0)
	data["beforeTime"] = query.GetParamOrDefault("before", 0)
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlCommentMusic, id), data, opt)
}
