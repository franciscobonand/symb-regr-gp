# Programação Genética - Regressão Simbólica

Regressão simbólica é uma técnica utilizada para descobrir fórmulas ou equações matemáticas que descrevem com certa precisão um conjunto de dados de entrada e saída.  

Programação genética é uma forma de computação evolutiva inspirada no processo de seleção natural.
Ela usa operadores genéticos, como mutação, crossover e seleção, para evoluir uma população de soluções candidatas em direção a uma solução ótima.  

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
Além da árvore, que representa seu genoma, a estrutura de um indivíduo também armazena o valor da Fitness, a altura da árvore e um campo booleano que diz se o indivíduo possui fitness válida ou não (a fitness só é inválida em casos de exceção ao criar/avaliar um indivíduo).   

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

- Grow: Os `Opcode`s da árvore são escolhidos considerando elementos dos conjuntos de variáveis e funções, considerando a altura máxima. Produz árvores com formas irregulares.
- Full: Os `Opcode`s da árvore são escolhidos apenas considerando as funções, até que a profundidade máxima da árvore seja alcançada. A partir desse momento, é considerado apenas o conjunto de variáveis. Produz árvores balanceadas.

### Fitness

A fitness é definida utilizando o método do erro quadrático médio ao se avaliar a `Expr` de um indivíduo, substituindo as variáveis pelos respectivos valores dados como entrada do programa.
Dessa forma, o valor da fitness sempre será um número não negativo, e quanto mais próxima de 0 (zero) for a fitness de um indivíduo, melhor ela é.

O erro quadrático médio é definido pela seguinte fórmula:  

![Fórmula RMSE](/images/rmse.png "Fórmula do erro quadrático médio")

Onde:
- *Ind* é o indivíduo sendo avaliado;
- EVAL(*Ind, x*) avalia a `Expr` do indivíduo *Ind* no conjunto de entrada fornecido *x*;
- *y* é a saída correta da função para a entrada *x*
- *N* é o número de exemplos fornecidos.

No caso da implementação feita, um indivíduo é avaliado com todos os dados fornecidos como entrada para o programa.

### Métodos de seleção

Nesse programa, foram implementados os métodos de seleção Aleatório, Roleta, Torneio e Lexicase.  
Todos os métodos implementados podem ser utilizados com ou sem elitismo, e a quantidade de indivíduos do elitismo é um dos parâmetros de execução do programa.

#### Aleatório
 
Esse método foi implementado para servir como base na análise dos demais métodos de seleção. Espera-se que ele possua o pior desempenho entre todos.  
Na seleção aleatória, indivíduos da população pai são randomicamente selecionados para comporem a próxima geração.

#### Roleta

No método da roleta, ou seleção proporcional à fitness, indivíduos são selecionados aleatoriamente da população pai, porém aqueles com melhor fitness têm maior probabilidade de serem escolhidos.

![Roleta](/images/rol-selection.svg "Seleção por Roleta")

**É válido ressaltar que o uso do erro quadrático médio para cálculo da fitness implica que indivíduos com menor fitness são melhores.**  

O cálculo das proporções da roleta é feito da seguinte forma:

- Obtem-se o valor da soma de todas as fitness de uma população (`fitSum`);
- A partir desse valor, obtem-se o valor da soma do percentual que a fitness de cada indivíduo representa da fitness total (`percSum`);
    - Como fitness menores são melhores, para cada indivíduo o cálculo feito é `percSum += (1 - fitnessDoIndividuo/fitSum)`;
- É escolhido um número aleatório entre 0.0 e `percSum` (`val`);
- Enquanto `val` for positivo, escolhe-se um indivíduo aleatório e é subtraído de `val` o valor do percentual que a fitness do indivíduo representa da fitness total;
    - `val -= (1 - fitnessDoIndividuo/fitSum)`;
- Quando essa subtração fizer com que `val` seja menor ou igual a zero, o indivíduo é adicionado à nova população. Isso é repetido até que a nova população tenha a quantidade de indivíduos igual à da população inicial;

#### Torneio

No método do torneio, primeiramente define-se o tamanho de torneio `K` (esse parâmetro é definido na execução do programa).  
Para cada iteração, são escolhidos aleatoriamente `K` indivíduos da população pai, e dentre esses indivíduos apenas o melhor deles é adicionado na nova população.  
Esse processo é repetido, com reposição, até a nova população tenha a mesma quantidade de indivíduos do que a população pai.

![Torneio](/images/tour-selection.svg "Seleção por Torneio")

É importante ressaltar que, quanto maior for o valor de `K`, maior a pressão seletiva é empregada, e a diversidade da população é diminuída.
Isso dá-se pelo fato de que, com um valor de `K` grande, há maior probabilidade de um mesmo melhor indivíduo ser selecionado a cada iteração do torneio.

#### Lexicase

No método lexicase, os indivíduos são avaliados com base em um subconjunto aleatório dos dados fornecidos como entrada para o problema,
e apenas indivíduos com a melhor fitness para esse subconjunto são adicionados à nova população.

A execução do lexicase que foi implementada nesse programa é feita da seguinte forma:

1. Inicialmente, todos os indivíduos da população são considerados candidatos para seleção;
2. É selecionado, de forma aleatória, um dos casos do conjunto de exemplos fornecidos como entrada para o programa;
3. Indivíduos candidatos são avaliados para esse caso, e aqueles com fitness pior que a melhor fitness para esse caso são removidos do conjunto de candidatos;
4. Se houver mais de um indivíduo no conjunto de candidatos, o caso atual é removido do conjunto de exemplos e os passos 2 e 3 são repetidos.
Se houver apenas um indivíduo no conjunto de candidatos, ele é adicionado à nova população.
Se não houverem mais exemplos a serem avaliados, escolhe-se um indivíduo aleatoriamente do conjunto de candidatos;

![Lexicase](/images/lex-selection.svg "Seleção Lexicase, com 1 indivíduo restante no conjunto de candidatos")

### Operadores genéticos

Os operadores genéticos presentes nessa implementação são Mutação e Crossover.
As probabilidades de indivíduos sofrerem mutação ou crossover são independentes, ou seja, um indivíduo pode passar tanto por crossover quanto por mutação em uma iteração.
Além disso, caso os novos indivíduos gerados a partir da aplicação dos operadores genéticos extrapolem o tamanho máximo permitido de seu genoma, eles são substituídos por seus pais.  

Optou-se por adicionar uma maior pressão seletiva, de forma que novos indivíduos frutos de crossover e/ou mutação que possuam fitness pior que a de seus pais sejam substituídos pelos pais.  
Isso foi feito pois inicialmente, quando essa pressão seletiva não havia sido implementada, os valores da fitness do pior indivíduo e da fitness média da população eram muito erráticos.
Isto é, após a evolução por diversas gerações, a fitness do melhor indivíduo convergia, mas as demais fitness possuíam variações absurdas de valor a cada geração 
(na geração *N-1* a fitness média poderia estar próxima da melhor fitness, mas na geração *N* mudar para um valor 1000x maior que a melhor fitness).

#### Mutação

Para realizar a mutação, primeiramente é selecionado de forma aleatória um nó da árvore do indivíduo.
Em seguida, é gerada uma nova árvore utilizando o método *Ramped half-and-half*.
Por fim, o nó que havia sido selecionado é substituído por essa nova árvore que foi gerada.

![Mutação](/images/mutation.svg "Mutação")

#### Crossover

Para realizar o crossover entre dois indivíduos, é selecionada uma subárvore aleatoriamente da árvore de cada um dos indivíduos.
Em seguida, a subárvore do indivíduo 2 é colocada no lugar da subárvore do indivíduo 1, e vice-versa.

![Crossover](/images/cx.svg "Crossover")

## Análises

As análises realizadas a partir dos [dados fornecidos](/datasets) podem ser encontradas no [Jupyter Notebook presente nesse repositório](CompNatTP1.ipynb).  
Alternativamente, as análises podem ser visualizadas acessando [esse link](https://colab.research.google.com/drive/1AWovkcjQS9xYW5QSmYPZGCevDP12N0un?usp=sharing).  
Os resultados usados para análise podem ser encontrados na pasta ['analysis'](/analysis).  

## Conclusão

O programa implementado, que usa programação genética para resolver o problema de regressão simbólica, apresentou resultados interessantes.
Ao combinar programação genética e regressão simbólica, foi capaz de descobrir automaticamente fórmulas matemáticas que descrevem com boa precisão um determinado conjunto de dados de entrada e saída,
principalmente dado que foram utilizadas apenas funções de adição, multiplicação, subtração e adição.

As experimentações com diferentes métodos de seleção e operadores genéticos revelaram algumas resultados interessantes.
A pressão seletiva adicionada nos operadores genéticos teve um impacto positivo no programa, fazendo com que a fitness geral da população convergisse.
Além disso, pode-se observar que a seleção lexicase começou a convergir para o resultado final mais rapidamente do que a seleção por torneio, mas, no final, o torneio teve um resultado melhor.

Ademais, o uso de elitismo nos métodos de seleção foi positivo, gerando resultados melhores para todos os métodos de seleção.
Por fim, também foi observado que uma maior chance de mutação dos indivíduos parece ter um impacto mais positivo no resultado final do que aumentar a probabilidade de crossover.

No geral, o programa fornece uma ferramenta boa e relativamente simples para descobrir fórmulas matemáticas que descrevem dados de entrada e saída.
Espera-se que esta documentação tenha fornecido uma visão abrangente do algoritmo do programa, seus recursos e como usá-lo para resolver o problema de regressão simbólica.

