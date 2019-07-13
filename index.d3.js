const antDirection = {
  up: 0,
  down: 1,
  left: 2,
  right: 3,
}


class LangtonAnt {
  constructor(rows, cols, rate) {
    this.rows = rows
    this.cols = cols
    this.ant = {x: parseInt(rows/2), y: parseInt(cols/2), v: 0, direction: antDirection.up}
    this.rate = rate  // living rate at the initBoard step
    this._lives = []
    this._end = false
  }

  initBoard = () => {
    this._end = false
    this._lives = _.range(this.rows)
      .map((r) => (
        _.range(this.cols).map((c) => Math.random() < this.rate ? 1 : 0)
      ))
    this.updateAntCell()
  }

  getLives = () => this._lives

  isStopped = () => this._end

  nextRound = () => {
    if (this._end) {
      return
    }
    // 若蚂蚁在白格，右转90度，将该格改为黑格，向前移一步；
    // 若蚂蚁在黑格，左转90度，将该格改为白格，向前移一步。
    if (this.ant.v === 0) {
      this.antTurnRight()
    } else {
      this.antTurnLeft()
    }
    this.updateAntCell()
  }

  updateAntCell = () => {
    if (this._end) {
      return
    }
    this.ant.v = this._lives[this.ant.x][this.ant.y]
    this._lives[this.ant.x][this.ant.y] = 2
  }

  antTurnRight = () => {
    this._lives[this.ant.x][this.ant.y] = 1
    switch (this.ant.direction) {
      case antDirection.up:
        this.ant.direction = antDirection.right
        this.ant.y++
        break
      case antDirection.right:
        this.ant.direction = antDirection.down
        this.ant.x++
        break
      case antDirection.down:
        this.ant.direction = antDirection.left
        this.ant.y--
        break
      case antDirection.left:
        this.ant.direction = antDirection.up
        this.ant.x--
        break
    }

    if (this.ant.x < 0 || this.ant.x >= this.rows || this.ant.y < 0 || this.ant.y >= this.cols) {
      this._end = true
    }
  }

  antTurnLeft = () => {
    this._lives[this.ant.x][this.ant.y] = 0
    switch (this.ant.direction) {
      case antDirection.up:
        this.ant.direction = antDirection.left
        this.ant.y--
        break
      case antDirection.right:
        this.ant.direction = antDirection.up
        this.ant.x--
        break
      case antDirection.down:
        this.ant.direction = antDirection.right
        this.ant.y++
        break
      case antDirection.left:
        this.ant.direction = antDirection.down
        this.ant.x++
        break
    }

    if (this.ant.x < 0 || this.ant.x >= this.rows || this.ant.y < 0 || this.ant.y >= this.cols) {
      this._end = true
    }
  }
}


class Board {
  constructor(container) {
    this.chart = {left: 20, top: 20, r: 4, border: 1}
    this.circleSize = this.chart.r * 2 + this.chart.border
    this.svg = d3.select(container).append('svg')
  }

  arrayToObject = (arr) => (
    _.flatten(_.range(arr.length).map((x) => (
      _.range(arr[0].length).map((y) => {
        return {x: x, y: y, v: arr[x][y]}
      })
    )))
  )

  render = (lives) => {
    const data = this.arrayToObject(lives)

    this.svg.attr('height', this.circleSize * lives.length + this.chart.top)
    // Enter
    this.svg.selectAll('circle')
      .data(data)
      .enter()
        .append('circle')
        .attr('cx', (d) => (
            d.y * this.circleSize + this.chart.left
          ))
        .attr('cy', (d) => (
            d.x * this.circleSize + this.chart.top
          ))
        .attr('r', this.chart.r)
        .attr('data', (d) => d.v)

    // Update
    this.svg.selectAll('circle')
      .data(data)
      .attr('data', (d) => d.v)

    // Exit
    this.svg.selectAll('circle')
      .data(data)
      .exit()
        .remove()
  }
}
