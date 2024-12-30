# 時相論理オレオレまとめ

$K = (W, R)$: Kripke構造

$\pi = w_0 w_1 w_2 \cdots$ where $w_i \in W, (w_i, w_{i + 1}) \in R$: path

$AP$: set of atomic formulae

## 相互再帰的な定義

$$\begin{align*}
StateFormula := & AtomicFormula \ | \\
& \neg StateFormula \ | \\
& StateFormula \vee StateFormula \ | \\
& \mathbf{E}\ PathFormula \ | \\
& \mathbf{A}\ PathFormula \\
\\
PathFormula := & \mathbf{X}\ StateFormula \ | \\
& \mathbf{F}\ StateFormula \ | \\
& \mathbf{G}\ StateFormula \ | \\
& StateFormula \ \mathbf{U} \ StateFormula
\end{align*}$$

意味論 (疑似コード)

```ocaml
type state_formula =
    | AtomicFormula of atomic_formula
    | Not of state_formula
    | Or of state_formula * state_formula
    | E of path_formula
    | A of path_formula

type path_formula =
    | X of state_formula
    | F of state_formula
    | G of state_formula
    | U of state_formula * state_formula

`|=` (w : world, p : atomic_formula) : bool = V(w, p)

`|=` (w : world, φ : state_formula) : bool =
    match φ with
    | AtomicFormula p -> w |= p
    | Not p -> not w |= p
    | Or (φ, ψ) -> w |= φ || w |= ψ
    | E φ -> exists π : path_formula, π(0) = w && π |= φ  (* Exists *)
    | A φ -> forall π : path_formula, π(0) = w -> π |= φ  (* All *)

`|=` (π : path, φ : path_formula) : bool =
    match φ with
    | X φ -> π(1) |= φ                                      (* neXt *)
    | F φ -> exists i, i >= 0 && π(i) |= φ                  (* Future *)
    | G φ -> forall i, i >= 0 -> π(i) |= φ                  (* Globally *)
    | U (φ, ψ) -> exists i >= 0,
                      π(i) |= ψ &&
                      (forall i, 0 <= j < i -> π(j) |= φ)  (* Until *)
```

パス論理式の意味は自然言語で書くと、$\pi$ に沿って実行が行われたときに:

- $\mathbf{X}\ \phi$ は、ワンステップ先で $\phi$ が成り立つこと
- $\mathbf{F}\ \phi$ は、未来のどこかで $\phi$ が成り立つこと
- $\mathbf{G}\ \phi$ は、常に $\phi$ が成り立ち続けること
- $\phi\ \mathbf{U}\ \psi$ は、未来のどこかの時点で $\psi$ が成立し、かつ、その少なくとも直前までは $\phi$ が成立し続けていること。

$\phi\ \mathbf{U}\ \psi$ の定義について、文章やコードでは、$\phi$ は、 $\psi$ が成り立つときには成り立っていなくてもいいように読める($j < i$) 。
しかし、意味論では、 $\psi$ が成り立っているときにも $\phi$ が成り立っている必要があるように読める ($j \le i$)。
他の文献、具体的には[この資料](https://www.cs.utexas.edu/~moore/acl2/seminar/2010.05-19-krug/slides.pdf)や[この資料](https://www.depts.ttu.edu/cs/research/documents/ksattlc.pdf)を見ても $j < i$ が正しそう。

## 相互再帰を使わない定義

状態論理式に出てくる $\mathbf{E}, \mathbf{A}$ と パス論理式に出てくる $\mathbf{X}, \mathbf{F}, \mathbf{G}, \mathbf{U}$ を予め組み合わせることで、相互再帰を行わずに状態論理式を定義できる。

$$\begin{align*}
StateFormula := & AtomicFormula \ | \\
& \neg StateFormula \ | \\
& StateFormula \vee StateFormula \ | \\
& \mathbf{EX}\ StateFormula \ | \\
& \mathbf{AX}\ StateFormula \ | \\
& \mathbf{EF}\ StateFormula \ | \\
& \mathbf{AF}\ StateFormula \ | \\
& \mathbf{EG}\ StateFormula \ | \\
& \mathbf{AG}\ StateFormula \ | \\
& \mathbf{E}\ (StateFormula\ \mathbf{U}\ StateFormula) \ | \\
& \mathbf{A}\ (StateFormula\ \mathbf{U}\ StateFormula)
\end{align*}$$

さらに、実はこれらは $\mathbf{EX}, \mathbf{EG}, \mathbf{EU}$ だけで書けるらしい(証明略)。

## 充足検査

$\pi$ は無限に長いものがありえて、それをどう扱うかを考える必要があるが、状態論理式にもはやパス論理式は明示的には出てこない。

**Def** 充足集合(satisfying set):  
$K = (W, R, I, AP, V)$: Kripke model, $\phi$: 状態論理式  
$$S(\phi) := \{w \in W | w \models \phi\}$$
を$\phi$ の充足集合という。

**Def** 充足: 
$K$ が $\phi$ を充足するとは、初期条件 $I$ が充足集合に含まれていることをいう。すなわち
$$I \subset S(\phi).$$
ただし、書籍では $I$ を集合ではなく単一の条件としているので、 $I \in S(\phi)$.

### EX論理式の検査
$\mathbf{EX}\phi$ は、ワンステップ先で $\phi$ が成り立つようなパスが存在することを表していたので、 $S(\mathbf{EX}\phi) = \{w \in W | \exists w' \in W. w'\in S(\phi) \wedge (w, w') \in R \}$ 

もうちょい真面目に書くと。

$\pi \models \mathbf{X}\phi$ は、 $\pi = w_0 w_1 \cdots$ としたとき、 $w_1 \models \phi$ を表す。また、 $w_0\models\mathbf{E}\Phi$ は、 $\exists \pi = w_0 \cdots, \pi\models\Phi$ を表す (ただし、$\pi$を有効なパスとする)。なので、 $w\models\mathbf{EX}\phi$ は、 $\exists\pi = ww_1\cdots, w_1\models\phi$ を表す。 $\pi$ が有効なパスであることから、これはすなわち、 $w\models\mathbf{EX}\phi = (\exists w_1 \in W, (w, w_1) \in R \wedge w_1 \models \phi)$ と同じ。$S(\phi) = \{w \in W | w \models \phi\}$ に注意すると、 
$$S(\mathbf{EX}\phi) = \{w \in W | \exists w' \in W, w' \in S(\phi) \wedge (w, w') \in R\}.$$

これくらいであれば $w$ のひとつ先を辿るだけで検査が出来る。

### EU論理式の検査
$\pi = w_0 w_1 w_2 \cdots$ とすると、 $\pi\models\phi\mathbf{U}\psi = (\exists i\in\mathbb{N}, w_i \in S(\psi) \wedge \forall j \leq i, w_j \in S(\phi))$ なので、 
$$S(\mathbf{E}(\phi\mathbf{U}\psi)) = \{w_0 \in W | \exists \pi = w_0 w_1 \cdots, \exists i\in\mathbb{N}, w_i \in S(\psi) \wedge \forall j < i, w_j \in S(\phi)\}.$$

パスは無限に長くなりうるが、終端の $S(\psi)$ を満たすものから検査していける。 (そのために、エッジを逆向きに辿るための情報を `kripkeModel` に追加する必要がある)  
すなわち、ある意味帰納的に求められる。

$$S(\mathbf{E}(\phi\mathbf{U}\psi)) = \{w | w \in S(\psi)\} \cup \{w \in S(\phi) | w' \in S(\mathbf{E}(\phi\mathbf{U}\psi)) \wedge (w, w') \in R\}$$

### EG論理式の検査
まず事前準備として、有向グラフの強連結分解アルゴリズムを用意する。

#### 有向グラフの強連結

有向グラフ $K = (W, R)$ が強連結であるとは、任意の $w, w' \in W$ がエッジをたどって到達可能なことをいう。  
強連結な部分グラフで極大のもの(すなわち、他のどの頂点を付け加えても強連結ではなくなるようなもの)を強連結成分という。
有向グラフは強連結成分で分解することができ、強連結成分分解という。また、強連結成分のうち、ノードを1つしか含まないものを自明な強連結成分という (自分自身にエッジがある場合も自明な強連結成分としていいのか疑問だが、現在取り扱っているKripke構造と意味論では、自分自身へのエッジは存在し得ないため考慮しない)

強連結成分分解を行うアルゴリズムに、例えば、Kosarajuのアルゴリズムがある。

**Kosarajuのアルゴリズム**

1. グラフを深さ優先探索→帰りがけ順を記録する
2. 帰りがけ順が後ろのものを出発点に、もとのグラフのエッジを逆向けにしたもので探索を行う。到達したものが強連結成分

強連結成分分解の応用などは [ここ](https://hcpc-hokudai.github.io/archive/graph_scc_001.pdf)が詳しい。

#### EG論理式の検査

分岐をうまく選択することで、目的の性質が常に成り立ち続けるような実行が可能、ということ。

$$S(\mathbf{EG}\phi) = \{w \in W | \exists \pi = w_0 w_1 \cdots, w_i \in S(\mathbf{EG}\phi)\ (w_i = w_0, w_1, \cdots)\}$$

パスは無限に長くなるかもしれないが、ノードは有限であるので、パスが訪ねるノードも有限。  
無限パスなら、ある時点から先は、同じルートを繰り返しぐるぐると回るはず。  
→ $S(\phi)$ に含まれる状態のみを経由し、 $S(\phi)$ の(非自明な)強連結成分にたどり着き、そのまま強連結成分の中でループに入る

このように考えると、EU論理式の検査と同様に、終端の、 $S(\phi)$ から逆向きに辿ることで探索ができる。