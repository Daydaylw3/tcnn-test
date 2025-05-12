#include <stdio.h>
#include <unistd.h>

void to_uppercase(char* str, int sleepTime) {
    sleep(sleepTime);
    for (int i = 0; str[i] != '\0'; i++) {
        if (str[i] >= 'a' && str[i] <= 'z') {
            str[i] = str[i] - ('a' - 'A');
        }
    }
}