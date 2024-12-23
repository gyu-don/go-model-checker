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
                      (forall i, 0 <= j <= i -> π(j) |= φ)  (* Until *)
```

パス論理式の意味は自然言語で書くと、$\pi$ に沿って実行が行われたときに:

- $\mathbf{X}\ \phi$ は、ワンステップ先で $\phi$ が成り立つこと
- $\mathbf{F}\ \phi$ は、未来のどこかで $\phi$ が成り立つこと
- $\mathbf{G}\ \phi$ は、常に $\phi$ が成り立ち続けること
- $\phi\ \mathbf{U}\ \psi$ は、未来のどこかの時点で $\psi$ が成立し、かつ、その少なくとも直前までは $\phi$ が成立し続けていること。

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