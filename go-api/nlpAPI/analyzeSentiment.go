
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

func GetTextSentiment(text string) ([]string, string, int, error) {
        
	ctx := context.TODO()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
			return nil, "", 0, err
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
		return nil, "", 0, err
	}

	type sentenceNLP struct {
		text string `json:"text"`
		score float64 `json:"score"`
	}

	fmt.Println(sentiment)
	fmt.Printf("Text: %v\n", text)


	// ネガティブな文章を集める
	var neg_phrase []string
	var minScore float32 = 100
	for _, sentence := range sentiment.Sentences {
		score := sentence.Sentiment.Score
		if score < -0.1 {
			neg_phrase = append(neg_phrase, sentence.Text.Content)
		}
		if minScore > score {
			minScore = score
		}
	}

	// scoreの数値を取得
	resScore := sentiment.DocumentSentiment.Magnitude * sentiment.DocumentSentiment.Score * 100
	
	// Return the judgment result
	switch {
	case minScore < -0.1:
		return neg_phrase, "neg", int(resScore), nil
	case minScore >= -0.1 && minScore <= 0.1:
		return make([]string, 0), "neutral", int(resScore), nil
	case minScore > 0.1:
		return make([]string, 0), "pos", int(resScore), nil
	}

	return neg_phrase, "neutral", int(resScore), nil

}