package models

type MyJsonName struct {
	Message struct {
		Tweet struct {
			Contributors interface{} `json:"contributors,omitempty"`
			Coordinates  interface{} `json:"coordinates omitempty"`
			CreatedAt    string      `json:"created_at omitempty"`
			Entities     struct {
				Hashtags     []interface{} `json:"hashtags omitempty"`
				Symbols      []interface{} `json:"symbols omitempty"`
				Urls         []interface{} `json:"urls omitempty"`
				UserMentions []interface{} `json:"user_mentions omitempty"`
			} `json:"entities"`
			FavoriteCount        int64       `json:"favorite_count,omitempty"`
			Favorited            bool        `json:"favorited,omitempty"`
			FilterLevel          string      `json:"filter_level,omitempty"`
			Geo                  interface{} `json:"geo,omitempty"`
			ID                   int64       `json:"id,omitempty"`
			IDStr                string      `json:"id_str,omitempty"`
			InReplyToScreenName  interface{} `json:"in_reply_to_screen_name,omitempty"`
			InReplyToStatusID    interface{} `json:"in_reply_to_status_id,omitempty"`
			InReplyToStatusIDStr interface{} `json:"in_reply_to_status_id_str,omitempty"`
			InReplyToUserID      interface{} `json:"in_reply_to_user_id,omitempty"`
			InReplyToUserIDStr   interface{} `json:"in_reply_to_user_id_str,omitempty"`
			IsQuoteStatus        bool        `json:"is_quote_status,omitempty"`
			Lang                 string      `json:"lang,omitempty"`
			Place                interface{} `json:"place,omitempty"`
			RetweetCount         int64       `json:"retweet_count,omitempty"`
			Retweeted            bool        `json:"retweeted,omitempty"`
			Source               string      `json:"source,omitempty"`
			Text                 string      `json:"text,omitempty"`
			TimestampMs          string      `json:"timestamp_ms,omitempty"`
			Truncated            bool        `json:"truncated,omitempty"`
			User                 struct {
				ContributorsEnabled            bool        `json:"contributors_enabled,omitempty"`
				CreatedAt                      string      `json:"created_at,omitempty"`
				DefaultProfile                 bool        `json:"default_profile,omitempty"`
				DefaultProfileImage            bool        `json:"default_profile_image,omitempty"`
				Description                    interface{} `json:"description,omitempty"`
				FavouritesCount                int64       `json:"favourites_count,omitempty"`
				FollowRequestSent              interface{} `json:"follow_request_sent,omitempty"`
				FollowersCount                 int64       `json:"followers_count,omitempty"`
				Following                      interface{} `json:"following,omitempty"`
				FriendsCount                   int64       `json:"friends_count,omitempty"`
				GeoEnabled                     bool        `json:"geo_enabled,omitempty"`
				ID                             int64       `json:"id,omitempty"`
				IDStr                          string      `json:"id_str,omitempty"`
				IsTranslator                   bool        `json:"is_translator,omitempty"`
				Lang                           string      `json:"lang,omitempty"`
				ListedCount                    int64       `json:"listed_count,omitempty"`
				Location                       string      `json:"location,omitempty"`
				Name                           string      `json:"name,omitempty"`
				Notifications                  interface{} `json:"notifications,omitempty"`
				ProfileBackgroundColor         string      `json:"profile_background_color,omitempty"`
				ProfileBackgroundImageURL      string      `json:"profile_background_image_url,omitempty"`
				ProfileBackgroundImageURLHTTPS string      `json:"profile_background_image_url_https,omitempty"`
				ProfileBackgroundTile          bool        `json:"profile_background_tile,omitempty"`
				ProfileImageURL                string      `json:"profile_image_url,omitempty"`
				ProfileImageURLHTTPS           string      `json:"profile_image_url_https,omitempty"`
				ProfileLinkColor               string      `json:"profile_link_color,omitempty"`
				ProfileSidebarBorderColor      string      `json:"profile_sidebar_border_color,omitempty"`
				ProfileSidebarFillColor        string      `json:"profile_sidebar_fill_color,omitempty"`
				ProfileTextColor               string      `json:"profile_text_color,omitempty"`
				ProfileUseBackgroundImage      bool        `json:"profile_use_background_image,omitempty"`
				Protected                      bool        `json:"protected,omitempty"`
				ScreenName                     string      `json:"screen_name,omitempty"`
				StatusesCount                  int64       `json:"statuses_count,omitempty"`
				TimeZone                       interface{} `json:"time_zone,omitempty"`
				URL                            interface{} `json:"url,omitempty"`
				UtcOffset                      interface{} `json:"utc_offset,omitempty"`
				Verified                       bool        `json:"verified,omitempty"`
			} `json:"user"`
		} `json:"tweet"`
		UnixTimestamp100us int64 `json:"unix_timestamp_100us,omitempty"`
	} `json:"message"`
}
