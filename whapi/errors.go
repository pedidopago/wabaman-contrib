package whapi

var (
	localizedErrorMap = map[int]string{
		0:      "Erro ao autenticar com o Meta. Contate o suporte.",
		3:      "Erro no token do telefone. Contate o suporte.",
		10:     "Erro na integração com o Meta. Contate o suporte.",
		190:    "O código de acesso expirou. Por favor, contate o suporte. (190)",
		4:      "Muitas tentativas de envio. Por favor, tente novamente em alguns minutos. (4)",
		80007:  "Erro ao enviar mensagem. Muitas tentativas de envio. Por favor, tente novamente em alguns minutos. (80007)",
		130429: "Erro ao enviar mensagem. Limite de conversas abertas pela empresa atingido. Por favor, tente novamente em alguns minutos. (130429)",
		131048: "Erro ao enviar mensagem. Limite de conversas abertas pela empresa atingido. (131048)",
		131056: "Erro ao enviar mensagem. Muitas tentativas de envio para este cliente. (131056)",
		133016: "Muitas tentativas de registro. (133016)",
		368:    "API temporariamente desativada devido a violações de políticas. Por favor, contate o suporte. (368)",
		131031: "Erro ao enviar mensagem. Muitas tentativas de envio. Por favor, tente novamente em alguns minutos. (131031)",
		1:      "Meta: Request Inválido ou possível erro de servidor. (1)",
		2:      "Meta: Serviços temporariamente indisponíveis. (2) https://metastatus.com/whatsapp-business-api",
		33:     "A conta deste telefone/branch foi excluída. Contate o suporte. (33)",
		100:    "Parâmetro inválido.",
		130472: "A mensagem não foi enviada pois um experimento do Meta está em andamento. Este cliente nunca receberá templates de marketing. (130472)",
		131000: "Erro interno da plataforma WhatsApp Business. Por favor, tente novamente mais tarde. (131000)",
		131005: "Acesso negado. Por favor, contate o suporte. (131005)",
		131008: "", //TODO: https://developers.facebook.com/docs/whatsapp/cloud-api/support/error-codes/
		131016: "Erro interno da plataforma WhatsApp Business. Por favor, tente novamente mais tarde. (131016)",
		131026: "Esta mensagem não pode ser entregue ao destinatário. Possíveis motivos: 1 - O número não está registrado no WhatsApp. 2 - O contato recebeu muitas mensagens de marketing de várias empresas em um curto período de tempo.",
		131042: "Erro de pagamento. Por favor, contate um gestor com acesso ao Meta Business. (131042)",
		131047: "Erro: Envia uma mensagem de reengajamento antes de enviar esta mensagem. (131047)",
		131049: "Erro: A Meta escolheu não entregar esta mensagem. (131049)",
		131052: "Erro ao baixar anexo do cliente. (131052)",
	}
)

func GetLocalizedError(code int) string {
	le, ok := localizedErrorMap[code]

	if !ok {
		return ""
	}

	return le
}
