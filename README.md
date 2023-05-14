# Programação Genética - Regressão Simbólica

Regressão simbólica é uma técnica utilizada para descobrir fórmulas ou equações matemáticas que descrevem com certa precisão um conjunto de dados de entrada e saída.  

Programação genética é uma forma de computação evolutiva inspirada no processo de seleção natural.
Ela usa operadores genéticos, como mutação, cruzamento e seleção, para evoluir uma população de soluções candidatas em direção a uma solução ótima.  

O programa aqui implementado combina essas duas técnicas para descobrir automaticamente fórmulas matemáticas que descrevem com precisão um determinado conjunto de dados de entrada e saída.
O programa começa gerando uma população inicial de indivíduos (que representam fórmulas aleatórias), que são evoluídas ao longo de múltiplas gerações usando operadores genéticos.
A aptidão (fitness) de cada fórmula é avaliada com base em quão bem ela se ajusta aos dados alvo.

Nesta documentação, é fornecida uma explicação detalhada do algoritmo do programa, seus recursos e como usá-lo para resolver o problema de regressão simbólica.
Além disso, é também disponibilizado um Jupyter Notebook contento análises após executar o programa com os dados de entrada que podem ser encontrados [nessa pasta](/datasets).
Os dados resultantes da execução do programa que foram utilizados nas análises podem ser encontrados [nessa pasta](/analysis).

## Como executar o programa

Primeiramente faça o download e instale a [versão mais recente da linguagem Golang](https://go.dev/doc/install).  
Com a instalação realizada, basta executar o seguinte comando do diretório raiz desse repositório:

```sh
go run .
```

**Caso não deseje instalar o Golang, pode optar por executar o binário que se encontra na pasta `/bin`**  
Para isso, primeiro execute o comando:

```sh 
chmod +x ./bin/symb-regr-gp
```

E então execute o programa com:

```sh 
./bin/symb-regr-gp
```

### Flags - Parametrização

Ao executar o programa, flags podem ser utilizadas para definir alguns parâmetros da execução.  
As flags são definidas da forma `.bin/symb-regr-gp -flag1 valor1 -flag2 valor2` (ou `go run . -flag1 valor1 -flag2 valor2`)  
São elas:

| Flag         | Default                          | Tipo            | Descrição                                               |
| ------------ | -------------------------------- | --------------- | ------------------------------------------------------- |
| \-popsize    | 20                               | Int > 0         | Tamanho da população                                    |
| \-gens       | 10                               | Int > 0         | Número de gerações a serem executadas                   |
| \-elitism    | 0                                | Int >= 0        | Número de indivíduos selecionados com elitismo          |
| \-selector   | tour                             | String          | Método de seleção ('rol', 'tour', 'lex' ou 'rand')      |
| \-toursize   | 2                                | Int >= 2        | Tamanho do Torneio (caso esse método seja usado)        |
| \-cxprob     | 0.9                              | 0 <= Float <= 1 | Probabilidade de realizar crossover                     |
| \-mutprob    | 0.05                             | 0 <= Float <= 1 | Probabilidade de realizar mutação                       |
| \-file       | datasets/synth1/synth1-train.csv | String          | Path para o arquivo de entrada do programa              |
| \-threads    | 1                                | Int > 0         | Quantidade de threads para avaliação em paralelo        |
| \-seed       | 1                                | Int             | Semente aleatória                                       |
| \-statsfile  | ""                               | String          | Gera relatório da execução e salva em arquivo informado |

Exemplo:

```sh
go run . -popsize 100 -selector tour -toursize 2 -gens 10 -threads 1 -file "datasets/synth1/synth1-train.csv" -cxprob 0.9 -mutprob 0.05 -elitism 1 -seed 1111
```

Também é possível ver a descrição das flags usando `--help`:

```sh
go run . --help
// OU
./bin/symb-regr-gp --help
```

