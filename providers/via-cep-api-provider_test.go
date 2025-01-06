// providers/via_cep_provider.go
package providers

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/rafael1abrao/goexpert/multithreading/domain"
)

func TestViaCepProvider_FetchCep(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		ctx context.Context
		cep string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.CepResponse
		wantErr bool
	}{
		{
			name: "Test ViaCepProvider FetchCep",
			fields: fields{
				client: http.DefaultClient,
			},
			args: args{
				ctx: context.Background(),
				cep: "01001000",
			},
			want: &domain.CepResponse{
				Cep:        "01001-000",
				Logradouro: "Praça da Sé",
				Bairro:     "Sé",
				Localidade: "São Paulo",
				UF:         "SP",
				Source:     "ViaCEP",
			},
			wantErr: false,
		},
		{
			name: "Test ViaCepProvider FetchCep with invalid CEP",
			fields: fields{
				client: http.DefaultClient,
			},
			args: args{
				ctx: context.Background(),
				cep: "000000000",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ViaCepProvider{
				client: tt.fields.client,
			}
			got, err := v.FetchCep(tt.args.ctx, tt.args.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViaCepProvider.FetchCep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ViaCepProvider.FetchCep() = %v, want %v", got, tt.want)
			}
		})
	}
}
