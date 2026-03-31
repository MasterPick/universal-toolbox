// Package daily 日常工具模块
// 提供计算器、单位换算、汇率换算、备忘录等日常实用工具
package daily

import (
	"fmt"
	"math"
	"strings"
	"time"
	"xtool/internal/db"
)

// DailyTools 日常工具结构体（Wails 绑定到前端）
type DailyTools struct {
	db *db.Database // 数据库连接
}

// Note 备忘录结构
type Note struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Color     string `json:"color"`
	Pinned    bool   `json:"pinned"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ConversionResult 单位转换结果
type ConversionResult struct {
	Value  float64 `json:"value"`  // 转换后的值
	Unit   string  `json:"unit"`   // 目标单位
	Formula string `json:"formula"` // 转换公式说明
}

// NewDailyTools 创建日常工具模块实例
func NewDailyTools(database *db.Database) *DailyTools {
	return &DailyTools{db: database}
}

// ============================================================
// 计算器（标准模式 + 科学模式）
// ============================================================

// CalcBasic 基础四则运算
// op: +, -, *, /
func (d *DailyTools) CalcBasic(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("除数不能为零")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("不支持的运算符: %s", op)
	}
}

// CalcScientific 科学计算函数
// func: sqrt/pow/log/ln/sin/cos/tan/abs/ceil/floor/round
func (d *DailyTools) CalcScientific(value, param float64, fn string) (float64, error) {
	switch fn {
	case "sqrt":
		if value < 0 {
			return 0, fmt.Errorf("负数不能开平方根")
		}
		return math.Sqrt(value), nil
	case "pow":
		return math.Pow(value, param), nil
	case "log":
		if value <= 0 {
			return 0, fmt.Errorf("对数的真数必须大于 0")
		}
		return math.Log10(value), nil
	case "ln":
		if value <= 0 {
			return 0, fmt.Errorf("自然对数的真数必须大于 0")
		}
		return math.Log(value), nil
	case "sin":
		return math.Sin(value * math.Pi / 180), nil // 角度转弧度
	case "cos":
		return math.Cos(value * math.Pi / 180), nil
	case "tan":
		return math.Tan(value * math.Pi / 180), nil
	case "abs":
		return math.Abs(value), nil
	case "ceil":
		return math.Ceil(value), nil
	case "floor":
		return math.Floor(value), nil
	case "round":
		return math.Round(value), nil
	case "pi":
		return math.Pi * value, nil
	default:
		return 0, fmt.Errorf("不支持的函数: %s", fn)
	}
}

// ============================================================
// 单位换算工具
// ============================================================

// ConvertLength 长度单位换算
// fromUnit/toUnit: mm, cm, m, km, inch, foot, yard, mile
func (d *DailyTools) ConvertLength(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为米（基准单位）
	toMeter := map[string]float64{
		"mm":   0.001,
		"cm":   0.01,
		"m":    1.0,
		"km":   1000.0,
		"inch": 0.0254,
		"foot": 0.3048,
		"yard": 0.9144,
		"mile": 1609.344,
	}

	fromFactor, ok1 := toMeter[fromUnit]
	toFactor, ok2 := toMeter[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的长度单位")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:  result,
		Unit:   toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertWeight 重量单位换算
// fromUnit/toUnit: mg, g, kg, t, oz, lb
func (d *DailyTools) ConvertWeight(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为克（基准单位）
	toGram := map[string]float64{
		"mg": 0.001,
		"g":  1.0,
		"kg": 1000.0,
		"t":  1000000.0,
		"oz": 28.3495,
		"lb": 453.592,
	}

	fromFactor, ok1 := toGram[fromUnit]
	toFactor, ok2 := toGram[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的重量单位")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:  result,
		Unit:   toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertTemperature 温度换算
// fromUnit/toUnit: C（摄氏度）、F（华氏度）、K（开尔文）
func (d *DailyTools) ConvertTemperature(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 先转换为摄氏度
	var celsius float64
	switch strings.ToUpper(fromUnit) {
	case "C":
		celsius = value
	case "F":
		celsius = (value - 32) * 5 / 9
	case "K":
		celsius = value - 273.15
	default:
		return ConversionResult{}, fmt.Errorf("不支持的温度单位（支持: C/F/K）")
	}

	// 再从摄氏度转为目标单位
	var result float64
	switch strings.ToUpper(toUnit) {
	case "C":
		result = celsius
	case "F":
		result = celsius*9/5 + 32
	case "K":
		result = celsius + 273.15
	default:
		return ConversionResult{}, fmt.Errorf("不支持的温度单位（支持: C/F/K）")
	}

	return ConversionResult{
		Value:  result,
		Unit:   toUnit,
		Formula: fmt.Sprintf("%g°%s = %g°%s", value, fromUnit, result, toUnit),
	}, nil
}

// ConvertSpeed 速度单位换算
// fromUnit/toUnit: ms（米/秒）、km/h（千米/时）、mph（英里/时）、knot（节）
func (d *DailyTools) ConvertSpeed(value float64, fromUnit, toUnit string) (ConversionResult, error) {
	// 统一转换为米/秒
	toMPS := map[string]float64{
		"ms":   1.0,
		"kmh":  1.0 / 3.6,
		"mph":  0.44704,
		"knot": 0.514444,
	}

	fromFactor, ok1 := toMPS[fromUnit]
	toFactor, ok2 := toMPS[toUnit]
	if !ok1 || !ok2 {
		return ConversionResult{}, fmt.Errorf("不支持的速度单位")
	}

	result := value * fromFactor / toFactor
	return ConversionResult{
		Value:  result,
		Unit:   toUnit,
		Formula: fmt.Sprintf("%g %s = %g %s", value, fromUnit, result, toUnit),
	}, nil
}

// ============================================================
// 时间日期工具
// ============================================================

// GetCurrentTime 获取当前时间的多种格式表示
func (d *DailyTools) GetCurrentTime() map[string]string {
	now := time.Now()
	return map[string]string{
		"datetime":  now.Format("2006-01-02 15:04:05"),
		"date":      now.Format("2006-01-02"),
		"time":      now.Format("15:04:05"),
		"weekday":   now.Weekday().String(),
		"timestamp": fmt.Sprintf("%d", now.Unix()),
		"utc":       now.UTC().Format("2006-01-02 15:04:05 UTC"),
		"timezone":  now.Format("MST"),
	}
}

// ============================================================
// 备忘录工具
// ============================================================

// GetNotes 获取所有备忘录
func (d *DailyTools) GetNotes() ([]Note, error) {
	rows, err := d.db.DB.Query(
		"SELECT id, title, content, color, pinned, created_at, updated_at FROM notes ORDER BY pinned DESC, updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		var pinnedInt int
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Color, &pinnedInt, &n.CreatedAt, &n.UpdatedAt); err != nil {
			continue
		}
		n.Pinned = pinnedInt == 1
		notes = append(notes, n)
	}

	return notes, nil
}

// SaveNote 保存备忘录（新增）
func (d *DailyTools) SaveNote(title, content, color string) (int64, error) {
	result, err := d.db.DB.Exec(
		"INSERT INTO notes (title, content, color) VALUES (?, ?, ?)",
		title, content, color,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateNote 更新备忘录内容
func (d *DailyTools) UpdateNote(id int64, title, content, color string) error {
	_, err := d.db.DB.Exec(
		"UPDATE notes SET title=?, content=?, color=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		title, content, color, id,
	)
	return err
}

// PinNote 切换备忘录置顶状态
func (d *DailyTools) PinNote(id int64) error {
	_, err := d.db.DB.Exec(
		"UPDATE notes SET pinned = CASE WHEN pinned = 1 THEN 0 ELSE 1 END WHERE id = ?",
		id,
	)
	return err
}

// DeleteNote 删除备忘录
func (d *DailyTools) DeleteNote(id int64) error {
	_, err := d.db.DB.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}
