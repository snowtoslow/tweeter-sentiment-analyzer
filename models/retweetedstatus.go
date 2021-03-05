package models

type RetweetedStatus struct {
	Contributors interface{} `json:"contributors"`
	Coordinates  interface{} `json:"coordinates"`
	CreatedAt    string      `json:"created_at"`
	Entities     struct {
		Hashtags     []interface{} `json:"hashtags"`
		Symbols      []interface{} `json:"symbols"`
		Urls         []interface{} `json:"urls"`
		UserMentions []interface{} `json:"user_mentions"`
	} `json:"entities"`
	FavoriteCount        int64       `json:"favorite_count"`
	Favorited            bool        `json:"favorited"`
	FilterLevel          string      `json:"filter_level"`
	Geo                  interface{} `json:"geo"`
	ID                   int64       `json:"id"`
	IDStr                string      `json:"id_str"`
	InReplyToScreenName  interface{} `json:"in_reply_to_screen_name"`
	InReplyToStatusID    interface{} `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr interface{} `json:"in_reply_to_status_id_str"`
	InReplyToUserID      interface{} `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   interface{} `json:"in_reply_to_user_id_str"`
	IsQuoteStatus        bool        `json:"is_quote_status"`
	Lang                 string      `json:"lang"`
	Place                interface{} `json:"place"`
	RetweetCount         int64       `json:"retweet_count"`
	Retweeted            bool        `json:"retweeted"`
	Source               string      `json:"source"`
	Text                 string      `json:"text"`
	Truncated            bool        `json:"truncated"`
	User                 User        `json:"user"`
	UniqueId             string      `json:"unique_id,omitempty"`
	AggregationRation    float64     `json:"aggregation_ration,omitempty"`
	SentimentSCore       int8        `json:"sentiment_s_core,omitempty"`
}
