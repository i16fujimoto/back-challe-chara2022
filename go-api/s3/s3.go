package s3

import (
	"os"
	"fmt"
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/joho/godotenv"
)

// S3のインスタンス作成
func newS3()(*s3.S3, error) {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	AWS_ACCESS_KEY := os.Getenv("MINIO_ACCESS_KEY")
	AWS_SECRET_KEY := os.Getenv("MINIO_SECRET_KEY")

	// sessionを作成
	// Must: NewSession関数を呼び出す際に、Sessionが有効であり、エラーがなかったことを確認するためのヘルパー関数
	// NewSession: SDKのデフォルト、設定ファイル、環境、ユーザー提供の設定ファイルから作成された新しいSession
	var newSession *session.Session = session.Must(session.NewSession())

	// aws(minio)のconfig情報の設定
	// 設定する値の型はawsのutility関数を利用して，ポインタに変換する必要がある
	config := aws.Config{
		Credentials: 					credentials.NewStaticCredentials(AWS_ACCESS_KEY, AWS_SECRET_KEY, ""), // 第3引数はtoken: STS経由でのセキュリティ証明書に必要，それ以外の場合は空文字
		Region:      					aws.String("ap-northest-1"),
		Endpoint: 	 					aws.String("http://minio:9000"), // 外部の場合は :9090
		S3ForcePathStyle: 				aws.Bool(true), // for Amazon S3: Virtual Hosting of Bucketsの場合のみ指定
		DisableRestProtocolURICleaning: aws.Bool(true), // URICleaningをDisableに設定：キーに隣接するスラッシュを含むオブジェクト（例：bucketname/foo/bar/objectname → bucketname + /foo/bar/objectnameに分割）を操作する場合のみ指定
	}

	// セッションを持つS3クライアントインスタンスを作成に返す
	return s3.New(newSession, &config), nil
}

// ダウンロードするオブジェクトの情報のゲッター
func getObjectInput(bucketName string, key string) *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
}

// S3にファイルをダウンロードし，その情報を返す
func download(s3_download *s3.S3, downloadKey *s3.GetObjectInput)([]byte, error) {

	/*
	（下記では，必ずオブジェクトをファイル化するときに利用する）
	
	f, err := os.Create("sample.txt")
	if err != nil {
		log.Fatal(err)
	}
	downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
        d.BufferProvider = s3manager.NewPooledBufferedWriterReadFromProvider(5 * 1024 * 1024) // ダウンロード処理の全体的なレイテンシを安定化
	})
	n, err := downloader.Download(f, &s3.GetObjectInput{ // 作成したfというファイルに書き込みを行う
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	*/
	
	image, err := s3_download.GetObject(downloadKey)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			// switch aerr.Code() {
			// case s3.ErrCodeNoSuchKey:
			// 	fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
			// case s3.ErrCodeInvalidObjectState:
			// 	fmt.Println(s3.ErrCodeInvalidObjectState, aerr.Error())
			// default:
			// 	fmt.Println(aerr.Error())
			// }
			return nil, aerr
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, err
		}
	}

	fmt.Println("S3からファイルのダウンロード: Done")

	//imageをbytes.Buffer型に変換します
	buf := new(bytes.Buffer)
	buf.ReadFrom(image.Body)

	return buf.Bytes(), nil

}


func getPutObjectInput(bucketName string, key string, imageFile *os.File) *s3.PutObjectInput {
	return &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:	aws.String(key),	
		Body: 	imageFile,
	}
}


// S3にファイルをアップロード
func upload(s3_upload *s3.S3, params *s3.PutObjectInput) error {

	_, err := s3_upload.PutObject(params)
	if err != nil {
		return err
	}

	fmt.Println("S3にアップロード: Done")

	return nil

}