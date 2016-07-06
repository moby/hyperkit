provider hyperkit {
	probe vmx__exit(int, unsigned int);
	probe vmx__ept__fault(int, unsigned long, unsigned long);
	probe vmx__inject__virq(int, int);
	probe vmx__write__msr(int, unsigned int, unsigned long);
	probe vmx__read__msr(int, unsigned int, unsigned long);
};
