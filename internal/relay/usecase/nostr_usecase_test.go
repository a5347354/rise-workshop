package usecase

import (
	"context"
	"testing"
)

func Test_relayUsecase_ReceiveMessage_Event(t *testing.T) {
	type args struct {
		ctx context.Context
		msg []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success Event",
			args: args{
				ctx: nil,
				msg: []byte(`["EVENT",{"id":"86d08bb580e0cc952c08dd37ba8abdbb93d92df8fcc19318eae3c2b505f7074f","pubkey":"ccd893700b5bad6b5d57dc35f7aca65b6ed3fcc8c8191616db41d1b61c953b91","created_at":1683885839,"kind":31337,"tags":[["d","b07v7s2ic0haospgmeg73i"],["media","https://media.zapstr.live:3118/d91191e30e00444b942c0e82cad470b32af171764c2275bee0bd99377efd4075/naddr1qqtxyvphwcmhxvnfvvcxsct0wdcxwmt9vumnx6gzyrv3ry0rpcqygju59s8g9jk5wzej4ut3wexzyad7uz7ejdm7l4q82qcyqqq856g4xnp7j","http"],["p","d91191e30e00444b942c0e82cad470b32af171764c2275bee0bd99377efd4075","Host"],["p","fa984bd7dbb282f07e16e7ae87b26a2a7b9b90b7246a44771f0cf5ae58018f52","Guest"],["c","Podcast"],["price","402"],["cover","https://s3-us-west-2.amazonaws.com/anchor-generated-image-bank/production/podcast_uploaded_nologo400/36291377/36291377-1673187804611-64b4f8e9f1687.jpg"],["subject","Nostrovia | The Pablo Episode"]],"content":"Nostrovia | The Pablo Episode\\n\\nhttps://s3-us-west-2.amazona","sig":"5977f56656760dc716af4b80510e2e53d5ae02fa25ec1f30c7747a0d11221739f0de3da4b3d24a595a19cc9e7236032163bfd8891faf73f7d1e1a2950bff006e"}]`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := relayUsecase{}
			if err := c.ReceiveMessage(tt.args.ctx, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("ReceiveMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
