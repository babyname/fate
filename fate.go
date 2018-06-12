package fate

type fate struct {
	name     *Name
	property *Property
}

func NewFate(lastName string) *fate {
	name := newName(lastName)
	return &fate{name: name}
}

func (f *fate) SetLastName(lastName string) {
	f.name = newName(lastName)
}

func (f *fate) SetLunarData() {

}

func (f *fate) SetProperty(p *Property) {
	//TODO
}

func (f *fate) MakeFirstName() *Name {

	return f.name
}
