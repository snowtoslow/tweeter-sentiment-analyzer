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
			FavoriteCount        int64           `json:"favorite_count,omitempty"`
			Favorited            bool            `json:"favorited,omitempty"`
			FilterLevel          string          `json:"filter_level,omitempty"`
			Geo                  interface{}     `json:"geo,omitempty"`
			ID                   int64           `json:"id,omitempty"`
			IDStr                string          `json:"id_str,omitempty"`
			InReplyToScreenName  interface{}     `json:"in_reply_to_screen_name,omitempty"`
			InReplyToStatusID    interface{}     `json:"in_reply_to_status_id,omitempty"`
			InReplyToStatusIDStr interface{}     `json:"in_reply_to_status_id_str,omitempty"`
			InReplyToUserID      interface{}     `json:"in_reply_to_user_id,omitempty"`
			InReplyToUserIDStr   interface{}     `json:"in_reply_to_user_id_str,omitempty"`
			IsQuoteStatus        bool            `json:"is_quote_status,omitempty"`
			Lang                 string          `json:"lang,omitempty"`
			Place                interface{}     `json:"place,omitempty"`
			RetweetCount         int64           `json:"retweet_count,omitempty"`
			Retweeted            bool            `json:"retweeted,omitempty"`
			RetweetedStatus      RetweetedStatus `json:"retweeted_status,omitempty"`
			Source               string          `json:"source,omitempty"`
			Text                 string          `json:"text,omitempty"`
			TimestampMs          string          `json:"timestamp_ms,omitempty"`
			Truncated            bool            `json:"truncated,omitempty"`
			User                 User            `json:"user" bson:"-"`
		} `json:"tweet"`
		UnixTimestamp100us int64   `json:"unix_timestamp_100us,omitempty"`
		UniqueId           string  `json:"unique_id,omitempty" bson:"-"`
		AggregationRation  float64 `json:"aggregation_ration,omitempty"`
		SentimentSCore     int8    `json:"sentiment_s_core,omitempty"`
		UserId             string  `json:"user_id,omitempty"`
	} `json:"message"`
}
