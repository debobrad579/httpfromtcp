OPERATIONS = ['+', '-', '×', '÷', '^', '%'];
TRIG_FUNCTIONS = ['sin', 'cos', 'tan'];
FUNCTIONS = [
    'log', 'ln', '√',  'asin', 'acos', 'atan', 'asinh', 'acosh', 'atanh', 
    'sin', 'cos', 'tan', 'sinh', 'cosh', 'tanh', 'exp'
];

class Calculator {
  constructor({ previousOperandOutput, currentOperandOutput, errorOutput, degButton, arcButton, hypButton }) {
    this.previousOperandOutput = previousOperandOutput;
    this.currentOperandOutput = currentOperandOutput;
    this.errorOutput = errorOutput;
    this.degButton = degButton;
    this.arcButton = arcButton;
    this.hypButton = hypButton;
    this.clear();
    this.setDeg(false)
  }

  clear() {
    this.currentOperand = '';
    this.previousOperand = '';
    this.justComputed = false;
    this.setArc(false)
    this.setHyp(false)
    this.updateDisplay();
  }

  delete() {
    if (this.justComputed) {this.clear(); return}

    for (let i = 0; i < FUNCTIONS.length; i++) {
      if (this.currentOperand.endsWith(FUNCTIONS[i] + '(')) {
        this.currentOperand = this.currentOperand.slice(0, -FUNCTIONS[i].length - 1);
        this.updateDisplay();
        return;
      }
    }

    this.currentOperand = this.currentOperand.slice(0, -1);
    this.updateDisplay();
  }

  appendString(string) {
    if (string === '.' && this.currentOperand.slice(this.getLastOperationIndex()).includes('.')) {return}

    if (this.justComputed) {
      this.previousOperand = this.currentOperand
      if (!OPERATIONS.includes(string)) {this.currentOperand = ''}
      this.justComputed = false;
    }

    if (TRIG_FUNCTIONS.includes(string)) {
      if (this.arc) {string = 'a' + string}
      if (this.hyperbolic) {string += 'h'}
    }

    if (FUNCTIONS.includes(string)) {
      string += '(';
    }

    this.currentOperand += string;
    this.updateDisplay();
  }

  compute() {
    this.appendMultiplication();
    const previousOperand = this.currentOperandOutput.innerText

    fetch('/api', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        equation: this.currentOperand
        .replaceAll('×', '*')
        .replaceAll('÷', '/')
        .replaceAll('π', 'pi')
        .replaceAll('√', 'sqrt')
        .replaceAll('log(', 'log10(')
        .replaceAll('ln(', 'log('),
        is_degree_mode: this.degreeMode,
      })
    })
    .then(res => {
      if (!res.ok) {
        this.errorOutput.innerText = res.statusText;
        return
      }

      return res.json();
    })
    .then(data => {
      if (data == null) return;
      this.previousOperand = previousOperand + ' =';
      this.justComputed = true;
      this.currentOperand = data.evaluated_value;
      this.updateDisplay();
    });
  }

  updateDisplay() {
    this.setArc(false)
    this.setHyp(false)
    this.previousOperandOutput.innerText = this.previousOperand;
    this.currentOperandOutput.innerText = this.currentOperand === '' ? '0' : this.currentOperand;
    this.errorOutput.innerText = '';
  }

  tokenize(expression) {
    const regex = /asinh|acosh|atanh|asin|acos|atan|sinh|cosh|tanh|sin|cos|tan|log|ln|√|exp|π|\d+\.\d+|\d+|[()+\-*/^×÷]/g;
    return expression.match(regex) || [];
  }

  appendMultiplication() {
    const input = this.currentOperand;
    const tokens = this.tokenize(input);
    const result = [];

    const isNumber = (t) => /^[0-9]+(\.[0-9]+)?$/.test(t) || t === 'π';

    for (let i = 0; i < tokens.length; i++) {
      const current = tokens[i];
      const next = tokens[i + 1];

      result.push(current);

      if (!next) continue;

      const insertMultiply =
        (
          (isNumber(current) || current === ')')
          &&
          (next === '(' || FUNCTIONS.includes(next) || isNumber(next))
        );

      if (insertMultiply) {
        result.push('*');
      }
    }

    this.currentOperand = result.join('');
  }

  getLastOperationIndex() {
    let lastOperationIndex = 0;

    for (let i = 0; i < OPERATIONS.length; i++) {
      const lastIIndex = this.currentOperand.lastIndexOf(OPERATIONS[i]);
      lastOperationIndex = lastIIndex > lastOperationIndex ? lastIIndex : lastOperationIndex;
    }

    return lastOperationIndex;
  }

  setArc(bool) {
    this.arc = bool;
    if (bool) {
      arcButton.classList = 'active';
    } else {
      arcButton.classList = '';
    }
  }

  setHyp(bool) {
    this.hyperbolic = bool;
    if (bool) {
      hypButton.classList = 'active';
    } else {
      hypButton.classList = '';
    }
  }

  setDeg(bool) {
    this.degreeMode = bool;
    if (bool) {
      degButton.innerText = "DEG";
    } else {
      degButton.innerText = "RAD";
    }
  }
}

const standardButtons = document.querySelectorAll('[data-button]');
const arcButton = document.querySelector('[data-arc]');
const hypButton = document.querySelector('[data-hyp]');
const degButton = document.querySelector('[data-deg]');
const equalsButton = document.querySelector('[data-equals]');
const deleteButton = document.querySelector('[data-delete]');
const allClearButton = document.querySelector('[data-all-clear]');
const previousOperandOutput = document.querySelector('[data-previous-operand]');
const currentOperandOutput = document.querySelector('[data-current-operand]');
const errorOutput = document.querySelector('[data-error]');

const calculator = new Calculator({ previousOperandOutput, currentOperandOutput, errorOutput, degButton, arcButton, hypButton });

for (let i = 0; i < standardButtons.length; i++) {
  const button = standardButtons[i];
  button.addEventListener('click', () => {calculator.appendString(button.innerText)})
}

arcButton.addEventListener('click', () => {calculator.setArc(!calculator.arc)})
hypButton.addEventListener('click', () => {calculator.setHyp(!calculator.hyperbolic)})
degButton.addEventListener('click', () => {calculator.setDeg(!calculator.degreeMode)})
equalsButton.addEventListener('click', () => {calculator.compute()})
deleteButton.addEventListener('click', () => {calculator.delete()})
allClearButton.addEventListener('click', () => {calculator.clear()})
