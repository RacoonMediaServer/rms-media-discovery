package media

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testDataSet1 = `
Unique ID : 248265037467786809018388670853630457298 (0xBAC61B6310A842678422AD6C1BD8F5D2)
Complete name : J:\Torrents\Готовые релизы\Матрица - Трилогия (1999-2003) WEBRip-AVC [Open Matte]\Матрица - Революция (2003) WEBRip-AVC [Open Matte]\Матрица - Революция (2003) WEBRip-AVC [Open Matte].mkv
Format : Matroska
Format version : Version 2
File size : 2.94 GiB
Duration : 2 h 9 min
Overall bit rate : 3 255 kb/s
Movie name : Матрица: Революция / The Matrix Revolutions (2003) [Open Matte]
Encoded date : UTC 2018-04-20 09:29:29
Writing application : mkvmerge v5.5.0 ('Healer') built on Apr 24 2012 23:47:03
Writing library : libebml v1.2.3 + libmatroska v1.3.0
Cover : Yes
Attachments : cover.pngVideo
ID : 1
Format : AVC
Format/Info : Advanced Video Codec
Format profile : High@L3.1
Format settings : CABAC / 8 Ref Frames
Format settings, CABAC : Yes
Format settings, ReFrames : 8 frames
Codec ID : V_MPEG4/ISO/AVC
Duration : 2 h 9 min
Bit rate : 2 100 kb/s
Width : 1 056 pixels
Height : 594 pixels
Display aspect ratio : 16:9
Frame rate mode : Constant
Frame rate : 23.976 (24000/1001) FPS
Color space : YUV
Chroma subsampling : 4:2:0
Bit depth : 8 bits
Scan type : Progressive
Bits/(Pixel*Frame) : 0.140
Stream size : 1.84 GiB (63%)
Title : WEBRip (AVC) / MPEG-4 AVC/H.264 / 2100 kbps / 1056x594 / 23.976 fps
Writing library : x264 core 146 r2538 121396c
Encoding settings : cabac=1 / ref=8 / deblock=1:0:0 / analyse=0x3:0x113 / me=umh / subme=10 / psy=1 / psy_rd=1.00:0.00 / mixed_ref=1 / me_range=32 / chroma_me=1 / trellis=2 / 8x8dct=1 / cqm=0 / deadzone=21,11 / fast_pskip=0 / chroma_qp_offset=-2 / threads=6 / lookahead_threads=1 / sliced_threads=0 / nr=0 / decimate=1 / interlaced=0 / bluray_compat=0 / constrained_intra=0 / bframes=9 / b_pyramid=2 / b_adapt=2 / b_bias=0 / direct=3 / weightb=1 / open_gop=0 / weightp=2 / keyint=250 / keyint_min=23 / scenecut=40 / intra_refresh=0 / rc_lookahead=50 / rc=2pass / mbtree=1 / bitrate=2100 / ratetol=1.0 / qcomp=0.60 / qpmin=0 / qpmax=69 / qpstep=4 / cplxblur=20.0 / qblur=0.5 / vbv_maxrate=14000 / vbv_bufsize=14000 / nal_hrd=none / filler=0 / ip_ratio=1.40 / aq=1:1.00
Language : English
Default : Yes
Forced : YesAudio #1
ID : 2
Format : AC-3
Format/Info : Audio Coding 3
Codec ID : A_AC3
Duration : 2 h 9 min
Bit rate mode : Constant
Bit rate : 384 kb/s
Channel(s) : 6 channels
Channel positions : Front: L C R, Side: L R, LFE
Sampling rate : 48.0 kHz
Frame rate : 31.250 FPS (1536 SPF)
Bit depth : 16 bits
Compression mode : Lossy
Stream size : 355 MiB (12%)
Title : DUB (Rus) / AC3 / 6 ch / 384 kbps / 48 kHz
Language : Russian
Service kind : Complete Main
Default : Yes
Forced : YesAudio #2
ID : 3
Format : AC-3
Format/Info : Audio Coding 3
Codec ID : A_AC3
Duration : 2 h 9 min
Bit rate mode : Constant
Bit rate : 384 kb/s
Channel(s) : 6 channels
Channel positions : Front: L C R, Side: L R, LFE
Sampling rate : 48.0 kHz
Frame rate : 31.250 FPS (1536 SPF)
Bit depth : 16 bits
Compression mode : Lossy
Stream size : 355 MiB (12%)
Title : AVO Гаврилов (Rus) / AC3 / 6 ch / 384 kbps / 48 kHz
Language : Russian
Service kind : Complete Main
Default : No
Forced : NoAudio #3
ID : 4
Format : AC-3
Format/Info : Audio Coding 3
Codec ID : A_AC3
Duration : 2 h 9 min
Bit rate mode : Constant
Bit rate : 384 kb/s
Channel(s) : 6 channels
Channel positions : Front: L C R, Side: L R, LFE
Sampling rate : 48.0 kHz
Frame rate : 31.250 FPS (1536 SPF)
Bit depth : 16 bits
Compression mode : Lossy
Stream size : 355 MiB (12%)
Title : Original (Eng) / AC3 / 6 ch / 384 kbps / 48 kHz
Language : English
Service kind : Complete Main
Default : No
Forced : NoText #1
ID : 5
Format : UTF-8
Codec ID : S_TEXT/UTF8
Codec ID/Info : UTF-8 Plain Text
Title : Full (Rus) / S_TEXT/UTF8
Language : Russian
Default : No
Forced : NoText #2
ID : 6
Format : UTF-8
Codec ID : S_TEXT/UTF8
Codec ID/Info : UTF-8 Plain Text
Title : Full (Eng) / S_TEXT/UTF8
Language : English
Default : No
Forced : NoText #3
ID : 7
Format : UTF-8
Codec ID : S_TEXT/UTF8
Codec ID/Info : UTF-8 Plain Text
Title : Full [SDH] (Rus) / S_TEXT/UTF8
Language : English
Default : No
Forced : NoMenu
00:00:00.000 : en:20 Hours to Go
00:03:55.402 : en:Trapped
00:07:33.119 : en:The Connection Matters
00:10:19.953 : en:Down Here I'm God
00:14:44.383 : en:Coat-Check Chaos
00:17:58.744 : en:Interesting Deal
00:23:33.245 : en:The End Is Coming
00:30:06.304 : en:We Meet at Last
00:35:35.633 : en:The Logos Found
00:39:34.205 : en:Volunteers
00:42:58.910 : en:I Believe in Him
00:46:43.300 : en:Nothing's Changed
00:49:02.606 : en:Stowaway
00:51:53.610 : en:See Your Enemy
00:58:00.644 : en:Give 'em Hell
01:00:19.115 : en:Detected
01:02:29.746 : en:Breaching the Dome
01:08:00.910 : en:Storming Sentinels
01:10:29.225 : en:Trying to Keep Up
01:14:33.803 : en:Unfinished Business
01:18:23.365 : en:Hell of a Pilot
01:23:08.984 : en:Report to the Council
01:25:55.817 : en:Believers
01:29:13.348 : en:Glimpse of the Sky
01:31:51.673 : en:Saying What Matters
01:36:51.472 : en:A Truce
01:41:55.276 : en:It Ends Tonight
01:44:53.120 : en:Urban Splash
01:48:06.647 : en:Neo's Choice
01:51:30.851 : en:Inevitable and Over
01:55:05.899 : en:Real Peace
01:57:31.211 : en:Freedom and Sunlight
02:00:12.372 : en:End Credits
`

func TestParseInfo(t *testing.T) {
	mi := ParseInfo(testDataSet1)
	expected := &Info{
		Format: "Matroska",
		Video: []VideoTrack{
			{
				Width:       1056,
				Height:      594,
				Codec:       "V_MPEG4/ISO/AVC",
				AspectRatio: "16:9",
			},
		},
		Audio: []AudioTrack{
			{
				Codec:    "A_AC3",
				Language: "Russian",
				Voice:    "DUB (Rus) / AC3 / 6 ch / 384 kbps / 48 kHz",
			},
			{
				Codec:    "A_AC3",
				Language: "Russian",
				Voice:    "AVO Гаврилов (Rus) / AC3 / 6 ch / 384 kbps / 48 kHz",
			},
			{
				Codec:    "A_AC3",
				Language: "English",
				Voice:    "Original (Eng) / AC3 / 6 ch / 384 kbps / 48 kHz",
			},
		},
		Subtitle: []SubtitleTrack{
			{
				Codec:    "S_TEXT/UTF8",
				Language: "Russian",
			},
			{
				Codec:    "S_TEXT/UTF8",
				Language: "English",
			},
			{
				Codec:    "S_TEXT/UTF8",
				Language: "English",
			},
		},
	}
	assert.Equal(t, expected, mi)
}
