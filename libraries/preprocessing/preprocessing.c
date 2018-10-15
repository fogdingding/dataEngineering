#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define _inpiutFileName "dataset/ettoday.rec"
#define _outputFileName "dataset/wordList.rec"
#define _utf8Length 3

char *deleteHeadWhitespace(char *);
int utf8Strlen(char *);
int isAsciiChar(char *);
int splitStringInFile(char *, int, int, FILE *);

int main() {
    int isContent = 0, i;
    FILE *fin = fopen(_inpiutFileName, "r");
    FILE *fout = fopen(_outputFileName, "w");
    char *inputBuffer = (char *)malloc(10000 * sizeof(char));

    while (fgets(inputBuffer, 10000, fin) != NULL) {
        if (inputBuffer[0] == '@' && inputBuffer[1] == 'B') {
            isContent = 1;
        } else if (isContent == 1) {
            int splitPos = 0;
            inputBuffer = deleteHeadWhitespace(inputBuffer);

            for (i = 0; i < strlen(inputBuffer); i++) {               
                if (memcmp(&inputBuffer[i], "，", _utf8Length) == 0) {
                    splitPos = splitStringInFile(inputBuffer, i, splitPos, fout);
                } else if (memcmp(&inputBuffer[i], "。", _utf8Length) == 0) {
                    splitPos = splitStringInFile(inputBuffer, i, splitPos, fout);
                } else if (memcmp(&inputBuffer[i], "；", _utf8Length) == 0) {
                    splitPos = splitStringInFile(inputBuffer, i, splitPos, fout);
                } else if (memcmp(&inputBuffer[i], "？", _utf8Length) == 0) {
                    splitPos = splitStringInFile(inputBuffer, i, splitPos, fout);
                } else if (memcmp(&inputBuffer[i], "！", _utf8Length) == 0) {
                    splitPos = splitStringInFile(inputBuffer, i, splitPos, fout);
                }
            }

            isContent = 0;
        }
    }

    free(inputBuffer);
    fclose(fin);
    fclose(fout);
    return 0;
}

char *deleteHeadWhitespace(char *str) {
    int i;
    for (i = 0; i < strlen(str); i++) {
        if (str[i] != ' ') {
            memcpy(str, &str[i], strlen(str) - i + 1);
            break;
        }
    }
    return str;
}

int splitStringInFile(char *str, int pos, int splitPos, FILE *fout) {
    // char *substr = (char *)malloc(1000 * sizeof(char));

    char substr[10000] = {'\0'};

    if (splitPos == 0) {
        memcpy(substr, &str[splitPos], pos - splitPos);
    } else {
        memcpy(substr, &str[splitPos + 1] + 2, pos - splitPos - 3);
    }

    memcpy(substr, deleteHeadWhitespace(substr), 10000);

    if (utf8Strlen(substr) >= 5 && !isAsciiChar(substr)) {
        fprintf(fout, "%s\n", substr);
    }

    // Error : could not free memory from malloc immediately
    // Cause high memory cost (about 22 gb)
    // free(substr);

    return pos;
}

int utf8Strlen(char *str) {
    int i, len = 0;
    for (i = 0; str[i]; i++) {
        if ((str[i] & 0xc0) != 0x80) {
            ++len;
        }
    }
    return len;
}

int isAsciiChar(char *str) {
    if (str[0] == 0x09 || str[0] == 0x0A || str[0] == 0x0D || (0x20 <= str[0] && str[0] <= 0x7E)) {
        return 1;
    } else {
        return 0;
    }
}