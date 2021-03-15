package models

type User struct {
	ContributorsEnabled            bool        `json:"contributors_enabled"`
	CreatedAt                      string      `json:"created_at"`
	DefaultProfile                 bool        `json:"default_profile"`
	DefaultProfileImage            bool        `json:"default_profile_image"`
	Description                    interface{} `json:"description"`
	FavouritesCount                int64       `json:"favourites_count"`
	FollowRequestSent              interface{} `json:"follow_request_sent"`
	FollowersCount                 int64       `json:"followers_count"`
	Following                      interface{} `json:"following"`
	FriendsCount                   int64       `json:"friends_count"`
	GeoEnabled                     bool        `json:"geo_enabled"`
	ID                             int64       `json:"id"`
	IDStr                          string      `json:"id_str"`
	IsTranslator                   bool        `json:"is_translator"`
	Lang                           string      `json:"lang"`
	ListedCount                    int64       `json:"listed_count"`
	Location                       interface{} `json:"location"`
	Name                           string      `json:"name"`
	Notifications                  interface{} `json:"notifications"`
	ProfileBackgroundColor         string      `json:"profile_background_color"`
	ProfileBackgroundImageURL      string      `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool        `json:"profile_background_tile"`
	ProfileBannerURL               string      `json:"profile_banner_url"`
	ProfileImageURL                string      `json:"profile_image_url"`
	ProfileImageURLHTTPS           string      `json:"profile_image_url_https"`
	ProfileLinkColor               string      `json:"profile_link_color"`
	ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string      `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool        `json:"profile_use_background_image"`
	Protected                      bool        `json:"protected"`
	ScreenName                     string      `json:"screen_name"`
	StatusesCount                  int64       `json:"statuses_count"`
	TimeZone                       interface{} `json:"time_zone"`
	URL                            interface{} `json:"url"`
	UtcOffset                      interface{} `json:"utc_offset"`
	Verified                       bool        `json:"verified"`
	UniqueId                       string      `json:"unique_id,omitempty" bson:"-"`
}
