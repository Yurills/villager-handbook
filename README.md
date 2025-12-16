# Villager Handbook (Social Deduction AI Engine)

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)
![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg?style=flat-square)
![Status](https://img.shields.io/badge/Status-Prototype-orange?style=flat-square)

**Villager Handbook** is a high-performance probabilistic reasoning engine designed to solve social deduction games (such as *Werewolf*, *Mafia*, or *Among Us*) using **Bayesian State Estimation** and **Information Theory**.

Unlike traditional bots that rely on static heuristics or "if-then" rule sets, this engine models the game as a **Partially Observable Markov Decision Process (POMDP)**. It maintains a dynamic belief state over thousands of parallel game configurations ("worlds") to calculate the exact probability of every player's role in real-time, adapting instantly to lies, contradictions, and complex social signaling.

---

## 🧠 Core Architecture

The system is built upon three fundamental computer science concepts:

### 1. The Multiverse (Hypothesis Space Generation)
At initialization, the engine employs **Recursive Backtracking** to generate the complete set of mathematically valid role assignments.
* **Constraint Satisfaction:** It filters permutations based on game rules (e.g., "Exactly 2 Wolves, 1 Seer").
* **Result:** A "Prior Distribution" representing the exhaustive hypothesis space before any game actions occur.

### 2. Bayesian Belief Update (State Estimation)
The engine updates the probability of each world based on observed interactions using Bayes' Theorem.
* **Likelihood Function ($P(E|H)$):** Each action (e.g., "Player A accuses Player B") is assigned a likelihood weight based on the actor's hypothetical role in that specific world.
    * *Wolf Behavior:* Modeled with weights for strategies like "Bussing" (sacrificing a partner) or "Pocketing" (defending an enemy).
    * *Villager Behavior:* Modeled with weights for truth-telling and error rates.
* **Posterior Update:** Worlds inconsistent with new "Hard Facts" (e.g., a role reveal) are eliminated ($P=0$). Worlds consistent with behavioral patterns are up-weighted.

### 3. Entropy-Minimization Search (Decision Making)
To determine the optimal move, the AI performs a **One-Step Expectimax Lookahead**:
1.  **Simulation:** Simulates every possible legal move (e.g., investigating any surviving player).
2.  **Projection:** Calculates the potential future belief states for every possible outcome of that move.
3.  **Optimization:** Selects the action that minimizes the **Expected Shannon Entropy** (Confusion) of the system, thereby maximizing **Information Gain**.

---

## 🚀 Key Features

* **Constraint Solving:** Instantly identifies logical paradoxes (e.g., the "Pigeonhole Principle" when the number of role claims exceeds the number of available roles).
* **Behavioral Modeling:** Distinguishes between "Standard Wolf Play" and advanced strategies like "Self-Bussing" or "Distancing."
* **In-Memory Processing:** Optimized for speed; calculates probability distributions across 5,000+ parallel worlds in microseconds.
* **Endgame Solver:** Uses exclusionary logic to solve "Mexican Standoff" scenarios where lack of information (silence) becomes a statistical indicator of guilt.

---

## 🛠 Installation

Ensure you have **Go 1.22+** installed.

```bash
# Clone the repository
git clone [https://github.com/Yurills/villager-handbook.git](https://github.com/Yurills/villager-handbook.git)

# Navigate to the directory
cd villager-handbook

# Install dependencies
go mod tidy
