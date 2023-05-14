# Programação Genética - Regressão Simbólica

Regressão simbólica é uma técnica utilizada para descobrir fórmulas ou equações matemáticas que descrevem com certa precisão um conjunto de dados de entrada e saída.  

Programação genética é uma forma de computação evolutiva inspirada no processo de seleção natural.
Ela usa operadores genéticos, como mutação, cruzamento e seleção, para evoluir uma população de soluções candidatas em direção a uma solução ótima.  

O programa aqui implementado combina essas duas técnicas para descobrir automaticamente fórmulas matemáticas que descrevem com precisão um determinado conjunto de dados de entrada e saída.
O programa começa gerando uma população inicial de indivíduos (que representam fórmulas aleatórias), que são evoluídas ao longo de múltiplas gerações usando operadores genéticos.
A aptidão (fitness) de cada fórmula é avaliada com base em quão bem ela se ajusta aos dados alvo.

Nesta documentação, é fornecida uma explicação detalhada do algoritmo do programa, seus recursos e como usá-lo para resolver o problema de regressão simbólica.  
Além disso, é também disponibilizado um Jupyter Notebook contendo análises após executar o programa com os dados de entrada que podem ser encontrados [nessa pasta](/datasets).  
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

| Flag           | Default                          | Tipo            | Descrição                                               |
| -------------- | -------------------------------- | --------------- | ------------------------------------------------------- |
| \-popsize      | 20                               | Int > 0         | Tamanho da população                                    |
| \-gens         | 10                               | Int > 0         | Número de gerações a serem executadas                   |
| \-elitism      | 0                                | Int >= 0        | Número de indivíduos selecionados com elitismo          |
| \-selector     | tour                             | String          | Método de seleção ('rol', 'tour', 'lex' ou 'rand')      |
| \-toursize     | 2                                | Int >= 2        | Tamanho do Torneio (caso esse método seja usado)        |
| \-cxprob       | 0.9                              | 0 <= Float <= 1 | Probabilidade de realizar crossover                     |
| \-mutprob      | 0.05                             | 0 <= Float <= 1 | Probabilidade de realizar mutação                       |
| \-file         | datasets/synth1/synth1-train.csv | String          | Path para o arquivo de entrada do programa              |
| \-threads      | 1                                | Int > 0         | Quantidade de threads para avaliação em paralelo        |
| \-seed         | 1                                | Int             | Semente aleatória                                       |
| \-statsfile    | `""`                             | String          | Gera relatório da execução e salva em arquivo informado |

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

## Implementação

Nesse tópico serão apresentadas as principais estruturas utilizadas no programa, assim como decisões de implementação e limitações.

### Indivíduo

Os indivíduos de uma população são representados por uma árvore, de altura máxima 7.
Além da árvore que representa seu genoma, a estrutura de um indivíduo também armazena o valor da Fitness, a altura da árvore e um campo booleano que diz se o indivíduo possui fitness válida ou não (a fitness só é inválida em casos de exceção ao criar/avaliar um indivíduo).   

```go
// Individual is a member of the population. Code represents its genome
type Individual struct {
	Code         operator.Expr
	Fitness      float64
	FitnessValid bool
	depth        int
}
```

`Expr`, abreviatura de "expressão", é a representação da árvore de um indivíduo. `Expr` é similar a um [heap](https://pt.wikipedia.org/wiki/Heap), no qual os itens do array são uma estrutura chamada `Opcode`.  

`Opcode`s são usados para representar funções (nós intermediários) ou variáveis (nós terminais, folhas da árvore).
O conjunto de funções escolhidas para o programa foram adição, subtração, divisão e multiplicação.
O número de variáveis disponíveis (x0, x1, ..., xN) varia conforme a quantidade de variáveis presentes no arquivo de entrada usado na execução do programa. 

![Representação de um indivíduo](/images/indiv-representation.svg "Heap, árvore derivada e expressão resultante")

#### Geração da população

Uma população, formada por um conjunto de indivíduos, é inicialmente gerada utilizando o método *Ramped half-and-half*.  
Esse método é basicamente uma combinação dos métodos *Grow* e *Full*.  

- Grow: O `Opcode` de um nó da árvore é escolhido considerando elementos dos conjuntos de variáveis e funções, considerando a altura máxima. Produz árvores com formas irregulares.
- Full: O `Opcode` de um nó da árvore é escolhido apenas considerando as funções, até que a profundidade máxima da árvore seja alcançada. A partir desse momento, é considerado apenas o conjunto de variáveis. Produz árvores balanceadas.


### Fitness



### Métodos de seleção



### Operadores genéticos

