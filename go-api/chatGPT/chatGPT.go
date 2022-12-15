package chatGPT

import (
	"os"
	"context"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

func Response(ctx context.Context, prompt []string) (string, error) {
	
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := gpt3.NewClient(apiKey)

	// CompletionWithEngine is the same as Completion except allows overriding the default engine on the client
	// Engine: "text-davinci-003"
	resp, err := client.CompletionWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: 			prompt, // Input
		MaxTokens: 			gpt3.IntPtr(512), // Output max tokens
		Temperature: 		gpt3.Float32Ptr(0.9), // 返答のランダム性の制御[0-1]
		N: 					gpt3.IntPtr(1), // 返答数
		Echo: 				false, // promptをエコーするか（返答のTextにpromptを含めるか）
		// Stop: []string{".", "。"},
		FrequencyPenalty: 	0.5, // 周波数制御[0-1] 高いと同じ話題を繰り返さなくなる（これまでのテキストに出現した頻度に応じてトークンにペナルティを課す）
		PresencePenalty: 	0.5,  // 新規トピック制御[0-1]：高いと新規のトピックが出現しやすくなる（これまでにテキストに出現したトークンにペナルティを課す）
	})

	if err != nil {
		return "", err
		// fmt.Println(gpt3.APIErrorResponse.Error)
	}

	/*
	type CompletionResponse struct {
		ID      string                     `json:"id"`
		Object  string                     `json:"object"`
		Created int                        `json:"created"`
		Model   string                     `json:"model"`
		Choices []CompletionResponseChoice `json:"choices"`
		Usage   CompletionResponseUsage    `json:"usage"`
	}

	type CompletionResponseChoice struct {
		Text         string        `json:"text"`
		Index        int           `json:"index"`
		LogProbs     LogprobResult `json:"logprobs"`
		FinishReason string        `json:"finish_reason"`
	}
	*/

	return resp.Choices[0].Text, nil
}