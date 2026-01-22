package contact

type Contact struct{}

func (c *Contact) HandlerName() string {
	return "contact"
}

func (c *Contact) DisplayName() string {
	return "Contacto"
}

func (c *Contact) RenderHTML() string {
	return `<!-- module -->
<section id="contact">
    <h1>Contáctanos</h1>
    <p>Envíanos un mensaje.</p>
</section>`
}

func (c *Contact) AllowedRoles(action byte) []byte {
	return []byte{'*'} // Public read
}

func (c *Contact) ValidateData(action byte, data ...any) error {
	return nil
}

func Add() []any {
	return []any{&Contact{}}
}
