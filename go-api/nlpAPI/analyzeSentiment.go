
package nlpAPI

import (
        "context"
        "fmt"
        
        language "cloud.google.com/go/language/apiv1"
        "cloud.google.com/go/language/apiv1/languagepb"
)

type SentenceNLP struct {
	text string `json:"text"`
	score float32 `json:"score"`
}

func GetTextSentiment(text string) ([]SentenceNLP, string, error) {
        
	ctx := context.TODO()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
			return nil, "", err
	}
	defer client.Close()

	// Detects the sentiment of the text.
	sentiment, err := client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
			Document: &languagepb.Document{
					Source: &languagepb.Document_Content{
							Content: text,
					},
					Type: languagepb.Document_PLAIN_TEXT,
			},
			EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		return nil, "", err
	}

	type sentenceNLP struct {
		text string `json:"text"`
		score float64 `json:"score"`
	}

	fmt.Println(sentiment)
	fmt.Printf("Text: %v\n", text)

// 	document_sentiment:{magnitude:1.7}  language:"ja"  sentences:{text:{content:"僕は優秀なエンジニア！"}  sentiment:{magnitude:0.9  score:0.9}}  sentences:{text:{content:"だけど，print文がわからないんだ"  begin_offset:33}  sentiment:{magnitude:0.7  score:-0.7}}

	// ネガティブな文章を集める
	var neg_phrase []SentenceNLP
	var minScore float32 = 100
	for _, sentence := range sentiment.Sentences {
		score := sentence.Sentiment.Score
		if score < -0.1 {
			neg_phrase = append(neg_phrase, SentenceNLP{
				text: sentence.Text.Content,
				score: score,
			})
		}
		if minScore > score {
			minScore = score
		}
	}
	
	// Return the judgment result
	switch {
	case minScore < -0.1:
		return neg_phrase, "negative", nil
	case minScore >= -0.1 && minScore <= 0.1:
		return nil, "neutral", nil
	case minScore > 0.1:
		return nil, "positive", nil
	}

	return neg_phrase, "neutral", nil

}