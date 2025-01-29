package fs

type TypeMap map[Type]string

// TypeInfo contains human-readable descriptions for supported file formats
var TypeInfo = TypeMap{
	ImageRaw:        "Unprocessed Sensor Data",
	ImageDNG:        "Adobe Digital Negative",
	ImageJPEG:       "Joint Photographic Experts Group (JPEG)",
	ImageJPEGXL:     "JPEG XL",
	ImageThumb:      "Thumbnail Image",
	ImagePNG:        "Portable Network Graphics",
	ImageGIF:        "Graphics Interchange Format",
	ImageTIFF:       "Tag Image File Format",
	ImagePSD:        "Adobe Photoshop",
	ImageBMP:        "Bitmap",
	ImageMPO:        "Stereoscopic JPEG (3D)",
	ImageAVIF:       "AV1 Image File Format",
	ImageAVIFS:      "AV1 Image Sequence",
	ImageHEIF:       "High Efficiency Image File Format",
	ImageHEIC:       "High Efficiency Image Container",
	ImageHEICS:      "HEIC Image Sequence",
	ImageWebP:       "Google WebP",
	VideoWebM:       "Google WebM",
	VideoMP2:        "MPEG 2 (H.262)",
	VideoAVC:        "Advanced Video Coding (H.264, MPEG-4 Part 10)",
	VideoHEVC:       "High Efficiency Video Coding (H.265)",
	VideoHEV1:       "High Efficiency Video Coding (HEVC) Bitstream",
	VideoVVC:        "Versatile Video Coding (H.266)",
	VideoEVC:        "Essential Video Coding (MPEG-5 Part 1)",
	VideoAV1:        "AOMedia Video 1",
	VideoMOV:        "Apple QuickTime",
	VideoMP4:        "Multimedia Container (MPEG-4 Part 14)",
	VideoM4V:        "Apple iTunes Multimedia Container",
	VideoMXF:        "Material Exchange Format",
	VideoAVI:        "Microsoft Audio Video Interleave",
	VideoASF:        "Advanced Systems Format ",
	VideoWMV:        "Windows Media",
	VideoDV:         "DV Video",
	Video3GP:        "Mobile Multimedia Container (3G)",
	Video3G2:        "Mobile Multimedia Container (CDMA2000)",
	VideoFlash:      "Adobe Flash",
	VideoMKV:        "Matroska Multimedia Container",
	VideoMPG:        "Moving Picture Experts Group (MPEG)",
	VideoMJPG:       "Motion JPEG",
	VideoAVCHD:      "Advanced Video Coding High Definition (AVCHD)",
	VideoBDAV:       "Blu-ray MPEG-2 Transport Stream",
	VideoOGV:        "Ogg Media (OGG)",
	VectorSVG:       "Scalable Vector Graphics",
	VectorAI:        "Adobe Illustrator",
	VectorPS:        "Adobe PostScript",
	VectorEPS:       "Encapsulated PostScript",
	SidecarXMP:      "Adobe Extensible Metadata Platform",
	SidecarAAE:      "Apple Image Edits XML",
	SidecarXML:      "Extensible Markup Language",
	SidecarJSON:     "Serialized JSON Data (Exiftool, Google Photos)",
	SidecarYAML:     "Serialized YAML Data (Config, Metadata)",
	SidecarText:     "Plain Text",
	SidecarInfo:     "Info Text",
	SidecarMarkdown: "Markdown Formatted Text",
	TypeUnknown:     "Other",
}
