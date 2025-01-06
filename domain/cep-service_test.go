// domain/cep_service.go
package domain

import (
	"context"
	"reflect"
	"testing"
)

type MockCepProvider struct {
	cep        string
	logradouro string
	bairro     string
	localidade string
	uf         string
	source     string
}

func (m *MockCepProvider) FetchCep(ctx context.Context, rawCep string) (*CepResponse, error) {
	return &CepResponse{
		Cep:        m.cep,
		Logradouro: m.logradouro,
		Bairro:     m.bairro,
		Localidade: m.localidade,
		UF:         m.uf,
		Source:     m.source,
	}, nil
}

func TestNormalizeAndValidateCep(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test NormalizeAndValidateCep with valid CEP",
			args: args{
				cep: "01001-000",
			},
			want:    "01001000",
			wantErr: false,
		},
		{
			name: "Test NormalizeAndValidateCep with invalid CEP",
			args: args{
				cep: "000000000",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeAndValidateCep(tt.args.cep)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeAndValidateCep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizeAndValidateCep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchFasterCep(t *testing.T) {
	type args struct {
		ctx       context.Context
		rawCep    string
		providers []CepProvider
	}
	tests := []struct {
		name    string
		args    args
		want    *CepResponse
		wantErr bool
	}{
		{
			name: "Test FetchFasterCep with valid CEP",
			args: args{
				ctx:    context.Background(),
				rawCep: "01001-000",
				providers: []CepProvider{
					&MockCepProvider{
						cep:        "01001000",
						logradouro: "Praça da Sé",
						bairro:     "Sé",
						localidade: "São Paulo",
						uf:         "SP",
						source:     "Mock",
					},
					&MockCepProvider{
						cep:        "01001000",
						logradouro: "Praça da Sé",
						bairro:     "Sé",
						localidade: "São Paulo",
						uf:         "SP",
						source:     "Mock",
					},
				},
			},
			want: &CepResponse{
				Cep:        "01001000",
				Logradouro: "Praça da Sé",
				Bairro:     "Sé",
				Localidade: "São Paulo",
				UF:         "SP",
				Source:     "Mock",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchFasterCep(tt.args.ctx, tt.args.rawCep, tt.args.providers)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchFasterCep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchFasterCep() = %v, want %v", got, tt.want)
			}
		})
	}
}
