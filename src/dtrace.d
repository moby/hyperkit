provider hyperkit {
	probe vmx__exit(int, unsigned int);
	probe vmx__ept__fault(int, unsigned long, unsigned long);
	probe vmx__inject__virq(int, int);
};
