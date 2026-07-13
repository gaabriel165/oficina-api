package notification

import (
	"fmt"
	"strings"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/resend/resend-go/v3"
)

type ResendNotifier struct {
	client *resend.Client
	from   string
}

func NewResendNotifier(apiKey, from string) *ResendNotifier {
	return &ResendNotifier{client: resend.NewClient(apiKey), from: from}
}

func (n *ResendNotifier) NotifyOrderStatusChanged(customer *entity.Customer, order *entity.ServiceOrder) error {
	content := statusContentFor(order.Status())
	greetingName := firstName(customer.Name())

	params := &resend.SendEmailRequest{
		From:    n.from,
		To:      []string{customer.Email()},
		Subject: fmt.Sprintf("Oficina Mecânica · OS %s — %s", shortID(order.ID()), content.label),
		Html:    buildStatusEmailHTML(greetingName, order, content),
		Text:    buildStatusEmailText(greetingName, order, content),
	}

	_, err := n.client.Emails.Send(params)
	return err
}

type statusContent struct {
	label   string
	message string
	color   string
}

var statusContentByStatus = map[valueobject.OrderStatus]statusContent{
	valueobject.OrderStatusReceived:        {"Recebida", "Recebemos a sua ordem de serviço e ela entrou na nossa fila de atendimento.", "#2563eb"},
	valueobject.OrderStatusInDiagnosis:     {"Em diagnóstico", "Nossa equipe está avaliando o seu veículo para identificar o que precisa ser feito.", "#7c3aed"},
	valueobject.OrderStatusWaitingApproval: {"Aguardando aprovação", "Seu orçamento está pronto! Assim que você aprovar, damos início ao serviço.", "#d97706"},
	valueobject.OrderStatusInExecution:     {"Em execução", "Boa notícia: o serviço no seu veículo já começou.", "#0891b2"},
	valueobject.OrderStatusFinished:        {"Finalizada", "O serviço foi concluído e o seu veículo já está pronto para retirada.", "#16a34a"},
	valueobject.OrderStatusDelivered:       {"Entregue", "Seu veículo foi entregue. Obrigado pela confiança!", "#16a34a"},
	valueobject.OrderStatusCancelled:       {"Cancelada", "Sua ordem de serviço foi cancelada. Em caso de dúvidas, fale com a nossa equipe.", "#dc2626"},
}

func statusContentFor(status valueobject.OrderStatus) statusContent {
	if content, ok := statusContentByStatus[status]; ok {
		return content
	}
	return statusContent{label: status.String(), message: "O status da sua ordem de serviço foi atualizado.", color: "#2563eb"}
}

func firstName(name string) string {
	fields := strings.Fields(name)
	if len(fields) == 0 {
		return "cliente"
	}
	return fields[0]
}

func shortID(id string) string {
	if len(id) < 8 {
		return id
	}
	return id[:8]
}

func formatBRL(value float64) string {
	return "R$ " + strings.Replace(fmt.Sprintf("%.2f", value), ".", ",", 1)
}

func buildStatusEmailHTML(name string, order *entity.ServiceOrder, content statusContent) string {
	amountRow := ""
	if order.TotalAmount() > 0 {
		amountRow = fmt.Sprintf(`
              <tr>
                <td style="padding:16px 20px;border-top:1px solid #e2e8f0;font-size:13px;color:#64748b;">Valor do orçamento</td>
                <td style="padding:16px 20px;border-top:1px solid #e2e8f0;font-size:15px;color:#0f172a;text-align:right;font-weight:bold;">%s</td>
              </tr>`, formatBRL(order.TotalAmount()))
	}

	return fmt.Sprintf(`<!doctype html>
<html lang="pt-BR">
<body style="margin:0;padding:0;background-color:#f1f5f9;font-family:Arial,Helvetica,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f1f5f9;padding:24px 0;">
    <tr>
      <td align="center">
        <table role="presentation" width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%%;background-color:#ffffff;border-radius:12px;overflow:hidden;border:1px solid #e2e8f0;">
          <tr>
            <td style="background-color:#0f172a;padding:24px 32px;">
              <span style="color:#ffffff;font-size:18px;font-weight:bold;letter-spacing:0.5px;">Oficina Mec&acirc;nica</span>
            </td>
          </tr>
          <tr>
            <td style="padding:32px;">
              <p style="margin:0 0 8px;font-size:16px;color:#0f172a;">Ol&aacute;, %s!</p>
              <p style="margin:0 0 24px;font-size:15px;line-height:1.6;color:#475569;">%s</p>
              <table role="presentation" cellpadding="0" cellspacing="0" style="margin:0 0 24px;">
                <tr>
                  <td style="background-color:%s;color:#ffffff;font-size:13px;font-weight:bold;padding:8px 18px;border-radius:999px;">%s</td>
                </tr>
              </table>
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="border:1px solid #e2e8f0;border-radius:8px;">
                <tr>
                  <td style="padding:16px 20px;font-size:13px;color:#64748b;">Ordem de servi&ccedil;o</td>
                  <td style="padding:16px 20px;font-size:13px;color:#0f172a;text-align:right;font-family:'Courier New',monospace;">%s</td>
                </tr>
                <tr>
                  <td style="padding:16px 20px;border-top:1px solid #e2e8f0;font-size:13px;color:#64748b;">Status atual</td>
                  <td style="padding:16px 20px;border-top:1px solid #e2e8f0;font-size:13px;color:#0f172a;text-align:right;font-weight:bold;">%s</td>
                </tr>%s
              </table>
              <p style="margin:24px 0 0;font-size:13px;line-height:1.6;color:#94a3b8;">Voc&ecirc; pode acompanhar o andamento da sua ordem de servi&ccedil;o a qualquer momento usando o n&uacute;mero acima.</p>
            </td>
          </tr>
          <tr>
            <td style="background-color:#f8fafc;padding:20px 32px;border-top:1px solid #e2e8f0;">
              <p style="margin:0;font-size:12px;color:#94a3b8;">Este &eacute; um e-mail autom&aacute;tico, por favor n&atilde;o responda. Em caso de d&uacute;vidas, entre em contato com a nossa equipe.</p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, name, content.message, content.color, content.label, order.ID(), content.label, amountRow)
}

func buildStatusEmailText(name string, order *entity.ServiceOrder, content statusContent) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Olá, %s!\n\n", name)
	fmt.Fprintf(&builder, "%s\n\n", content.message)
	fmt.Fprintf(&builder, "Ordem de serviço: %s\n", order.ID())
	fmt.Fprintf(&builder, "Status atual: %s\n", content.label)
	if order.TotalAmount() > 0 {
		fmt.Fprintf(&builder, "Valor do orçamento: %s\n", formatBRL(order.TotalAmount()))
	}
	builder.WriteString("\nVocê pode acompanhar o andamento da sua ordem de serviço a qualquer momento usando o número acima.\n\n")
	builder.WriteString("Oficina Mecânica — este é um e-mail automático.")
	return builder.String()
}
