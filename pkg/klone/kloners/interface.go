package kloners

type Kloner interface {
	SetContext(context KlonerContext)
	Klone() (error)
}

type KlonerContext interface {

}
