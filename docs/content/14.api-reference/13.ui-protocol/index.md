---
title: UI Protocol API
description: Nomad UI Protocol 完整 API 参考
navigation:
  icon: i-lucide-code
---

# UI Protocol API

本文档提供 Nomad UI Protocol 的完整 API 参考。

## Go 类型定义

### NomadUIMessage

主消息结构，包含四种操作类型：

```go
type NomadUIMessage struct {
    SurfaceUpdate   *SurfaceUpdateMessage   `json:"surfaceUpdate,omitempty"`
    DataModelUpdate *DataModelUpdateMessage `json:"dataModelUpdate,omitempty"`
    BeginRendering  *BeginRenderingMessage  `json:"beginRendering,omitempty"`
    DeleteSurface   *DeleteSurfaceMessage   `json:"deleteSurface,omitempty"`
}
```

### SurfaceUpdateMessage

更新 Surface 的组件定义：

```go
type SurfaceUpdateMessage struct {
    // Surface 唯一标识符
    SurfaceID string `json:"surfaceId"`
    
    // 组件定义列表（增量合并）
    Components []ComponentDefinition `json:"components,omitempty"`
}
```

### DataModelUpdateMessage

更新数据模型：

```go
type DataModelUpdateMessage struct {
    // Surface 唯一标识符
    SurfaceID string `json:"surfaceId"`
    
    // JSON Pointer 路径，默认 "/"
    Path string `json:"path,omitempty"`
    
    // 数据内容
    Contents any `json:"contents"`
}
```

### BeginRenderingMessage

开始渲染 Surface：

```go
type BeginRenderingMessage struct {
    // Surface 唯一标识符
    SurfaceID string `json:"surfaceId"`
    
    // 根组件 ID
    Root string `json:"root"`
    
    // CSS 自定义属性（主题化）
    Styles map[string]string `json:"styles,omitempty"`
}
```

### DeleteSurfaceMessage

删除 Surface：

```go
type DeleteSurfaceMessage struct {
    // Surface 唯一标识符
    SurfaceID string `json:"surfaceId"`
}
```

### ComponentDefinition

组件定义（邻接表模型）：

```go
type ComponentDefinition struct {
    // 组件唯一标识符
    ID string `json:"id"`
    
    // 流式渲染权重："initial" | "final"
    Weight string `json:"weight,omitempty"`
    
    // 组件规格
    Component ComponentSpec `json:"component"`
}
```

### ComponentSpec

组件规格（联合类型）：

```go
type ComponentSpec struct {
    // 布局组件
    Row     *RowProps     `json:"Row,omitempty"`
    Column  *ColumnProps  `json:"Column,omitempty"`
    Card    *CardProps    `json:"Card,omitempty"`
    List    *ListProps    `json:"List,omitempty"`
    Tabs    *TabsProps    `json:"Tabs,omitempty"`
    Modal   *ModalProps   `json:"Modal,omitempty"`
    Divider *DividerProps `json:"Divider,omitempty"`
    
    // 内容组件
    Text        *TextProps        `json:"Text,omitempty"`
    Image       *ImageProps       `json:"Image,omitempty"`
    Icon        *IconProps        `json:"Icon,omitempty"`
    Video       *VideoProps       `json:"Video,omitempty"`
    AudioPlayer *AudioPlayerProps `json:"AudioPlayer,omitempty"`
    
    // 输入组件
    Button         *ButtonProps         `json:"Button,omitempty"`
    TextField      *TextFieldProps      `json:"TextField,omitempty"`
    Checkbox       *CheckboxProps       `json:"Checkbox,omitempty"`
    Select         *SelectProps         `json:"Select,omitempty"`
    DateTimeInput  *DateTimeInputProps  `json:"DateTimeInput,omitempty"`
    Slider         *SliderProps         `json:"Slider,omitempty"`
    MultipleChoice *MultipleChoiceProps `json:"MultipleChoice,omitempty"`
    
    // 自定义组件
    Custom *CustomProps `json:"Custom,omitempty"`
}
```

### PropertyValue

属性值类型（支持字面值和数据绑定）：

```go
type PropertyValue struct {
    // 字符串字面值
    LiteralString *string `json:"literalString,omitempty"`
    
    // 数字字面值
    LiteralNumber *float64 `json:"literalNumber,omitempty"`
    
    // 布尔字面值
    LiteralBoolean *bool `json:"literalBoolean,omitempty"`
    
    // JSON Pointer 数据绑定路径
    Path *string `json:"path,omitempty"`
}
```

### ComponentArrayReference

子组件引用：

```go
type ComponentArrayReference struct {
    // 显式 ID 列表
    ExplicitList []string `json:"explicitList,omitempty"`
    
    // 模板渲染
    Template *TemplateReference `json:"template,omitempty"`
}

type TemplateReference struct {
    // 模板组件 ID
    ComponentID string `json:"componentId"`
    
    // 数据源路径
    DataBinding string `json:"dataBinding"`
}
```

## 组件 Props

### 布局组件

#### RowProps

```go
type RowProps struct {
    Children ComponentArrayReference `json:"children"`
    Gap      int                     `json:"gap,omitempty"`
    Align    string                  `json:"align,omitempty"` // "start" | "center" | "end" | "stretch"
}
```

#### ColumnProps

```go
type ColumnProps struct {
    Children ComponentArrayReference `json:"children"`
    Gap      int                     `json:"gap,omitempty"`
    Align    string                  `json:"align,omitempty"` // "start" | "center" | "end" | "stretch"
}
```

#### CardProps

```go
type CardProps struct {
    Children ComponentArrayReference `json:"children"`
    Title    *PropertyValue          `json:"title,omitempty"`
    Subtitle *PropertyValue          `json:"subtitle,omitempty"`
}
```

#### ListProps

```go
type ListProps struct {
    Children ComponentArrayReference `json:"children"`
    Dividers bool                    `json:"dividers,omitempty"`
}
```

#### TabsProps

```go
type TabsProps struct {
    ActiveTab PropertyValue   `json:"activeTab"` // 必须是 path 类型
    Tabs      []TabDefinition `json:"tabs"`
}

type TabDefinition struct {
    ID      string                  `json:"id"`
    Label   PropertyValue           `json:"label"`
    Content ComponentArrayReference `json:"content"`
}
```

#### ModalProps

```go
type ModalProps struct {
    Open     PropertyValue           `json:"open"` // 必须是 path 类型
    Title    *PropertyValue          `json:"title,omitempty"`
    Children ComponentArrayReference `json:"children"`
}
```

#### DividerProps

```go
type DividerProps struct {
    Orientation string `json:"orientation,omitempty"` // "horizontal" | "vertical"
}
```

### 内容组件

#### TextProps

```go
type TextProps struct {
    Text      PropertyValue `json:"text"`
    UsageHint string        `json:"usageHint,omitempty"` // "h1" | "h2" | "h3" | "h4" | "h5" | "caption" | "body"
}
```

#### ImageProps

```go
type ImageProps struct {
    Src       PropertyValue `json:"src"`
    Alt       *PropertyValue `json:"alt,omitempty"`
    UsageHint string        `json:"usageHint,omitempty"` // "icon" | "avatar" | "smallFeature" | "mediumFeature" | "largeFeature" | "header"
}
```

#### IconProps

```go
type IconProps struct {
    Name  PropertyValue `json:"name"`
    Size  int           `json:"size,omitempty"`
    Color *PropertyValue `json:"color,omitempty"`
}
```

#### VideoProps

```go
type VideoProps struct {
    Src      PropertyValue  `json:"src"`
    Poster   *PropertyValue `json:"poster,omitempty"`
    Autoplay bool           `json:"autoplay,omitempty"`
    Controls bool           `json:"controls,omitempty"`
}
```

#### AudioPlayerProps

```go
type AudioPlayerProps struct {
    Src      PropertyValue `json:"src"`
    Title    *PropertyValue `json:"title,omitempty"`
    Controls bool          `json:"controls,omitempty"`
}
```

### 输入组件

#### ButtonProps

```go
type ButtonProps struct {
    Label    PropertyValue  `json:"label"`
    Action   string         `json:"action"`  // 动作标识符
    Variant  string         `json:"variant,omitempty"` // "primary" | "secondary" | "text"
    Disabled *PropertyValue `json:"disabled,omitempty"`
}
```

#### TextFieldProps

```go
type TextFieldProps struct {
    Value       PropertyValue  `json:"value"` // 必须是 path 类型
    Label       *PropertyValue `json:"label,omitempty"`
    Placeholder *PropertyValue `json:"placeholder,omitempty"`
    Multiline   bool           `json:"multiline,omitempty"`
    Disabled    *PropertyValue `json:"disabled,omitempty"`
}
```

#### CheckboxProps

```go
type CheckboxProps struct {
    Checked  PropertyValue  `json:"checked"` // 必须是 path 类型
    Label    *PropertyValue `json:"label,omitempty"`
    Disabled *PropertyValue `json:"disabled,omitempty"`
}
```

#### SelectProps

```go
type SelectProps struct {
    Value    PropertyValue  `json:"value"`   // 必须是 path 类型
    Options  PropertyValue  `json:"options"` // 选项数组路径
    Label    *PropertyValue `json:"label,omitempty"`
    Disabled *PropertyValue `json:"disabled,omitempty"`
}
```

#### DateTimeInputProps

```go
type DateTimeInputProps struct {
    Value    PropertyValue  `json:"value"` // 必须是 path 类型
    Label    *PropertyValue `json:"label,omitempty"`
    Type     string         `json:"type,omitempty"` // "date" | "time" | "datetime"
    Disabled *PropertyValue `json:"disabled,omitempty"`
}
```

#### SliderProps

```go
type SliderProps struct {
    Value    PropertyValue  `json:"value"` // 必须是 path 类型
    Label    *PropertyValue `json:"label,omitempty"`
    Min      float64        `json:"min,omitempty"`
    Max      float64        `json:"max,omitempty"`
    Step     float64        `json:"step,omitempty"`
    Disabled *PropertyValue `json:"disabled,omitempty"`
}
```

#### MultipleChoiceProps

```go
type MultipleChoiceProps struct {
    Value    PropertyValue  `json:"value"`   // 必须是 path 类型
    Options  PropertyValue  `json:"options"` // 选项数组路径
    Label    *PropertyValue `json:"label,omitempty"`
    Multiple bool           `json:"multiple,omitempty"`
    Disabled *PropertyValue `json:"disabled,omitempty"`
}
```

#### CustomProps

```go
type CustomProps struct {
    Type  string                   `json:"type"`  // 已注册的自定义组件类型
    Props map[string]PropertyValue `json:"props"`
}
```

## 事件类型

### ProgressUISurfaceUpdateEvent

```go
type ProgressUISurfaceUpdateEvent struct {
    SurfaceID  string                `json:"surface_id"`
    Components []ComponentDefinition `json:"components,omitempty"`
    Root       string                `json:"root,omitempty"`
    Styles     map[string]string     `json:"styles,omitempty"`
}

func (e *ProgressUISurfaceUpdateEvent) Channel() AgentChannel { return ChannelProgress }
func (e *ProgressUISurfaceUpdateEvent) EventType() string     { return "ui:surface_update" }
```

### ProgressUIDataUpdateEvent

```go
type ProgressUIDataUpdateEvent struct {
    SurfaceID string `json:"surface_id"`
    Path      string `json:"path"`
    Contents  any    `json:"contents"`
}

func (e *ProgressUIDataUpdateEvent) Channel() AgentChannel { return ChannelProgress }
func (e *ProgressUIDataUpdateEvent) EventType() string     { return "ui:data_update" }
```

### ProgressUIDeleteSurfaceEvent

```go
type ProgressUIDeleteSurfaceEvent struct {
    SurfaceID string `json:"surface_id"`
}

func (e *ProgressUIDeleteSurfaceEvent) Channel() AgentChannel { return ChannelProgress }
func (e *ProgressUIDeleteSurfaceEvent) EventType() string     { return "ui:delete_surface" }
```

### ControlUIActionEvent

```go
type ControlUIActionEvent struct {
    SurfaceID   string         `json:"surface_id"`
    ComponentID string         `json:"component_id"`
    Action      string         `json:"action"`
    Payload     map[string]any `json:"payload,omitempty"`
}

func (e *ControlUIActionEvent) Channel() AgentChannel { return ChannelControl }
func (e *ControlUIActionEvent) EventType() string     { return "ui:action" }
```

## TypeScript 类型

### MessageProcessor

```typescript
interface MessageProcessor {
  // 获取所有 surfaces
  getSurfaces(): ReadonlyMap<string, Surface>
  
  // 获取指定 surface
  getSurface(surfaceId: string): Surface | undefined
  
  // 清空所有 surfaces
  clearSurfaces(): void
  
  // 处理消息批次
  processMessages(messages: NomadUIMessage[]): void
  
  // 处理单个消息
  processMessage(message: NomadUIMessage): void
  
  // 获取数据
  getData(surfaceId: string, path: string): DataValue | null
  
  // 设置数据
  setData(surfaceId: string, path: string, value: DataValue): boolean
  
  // 订阅 surface 变化
  subscribe(surfaceId: string, listener: (surface: Surface) => void): () => void
  
  // 订阅 surface 删除
  subscribeToDelete(surfaceId: string, listener: (surfaceId: string) => void): () => void
  
  // 检查是否在流式模式
  isStreaming(surfaceId: string): boolean
}
```

### ComponentRegistry

```typescript
interface ComponentRegistry {
  // 注册组件
  register(typeName: string, constructor: ComponentConstructor, tagName?: string): void
  
  // 获取组件
  get(typeName: string): ComponentConstructor | undefined
  
  // 检查组件是否已注册
  has(typeName: string): boolean
  
  // 冻结注册表（生产模式）
  freeze(): void
  
  // 获取所有已注册的组件类型
  getRegisteredTypes(): string[]
}
```

### Surface

```typescript
interface Surface {
  rootComponentId: string | null
  componentTree: AnyComponentNode | null
  dataModel: DataMap
  components: Map<string, ComponentDefinition>
  styles: Record<string, string>
}
```

### UIActionEvent

```typescript
interface UIActionEvent {
  surfaceId: string
  componentId: string
  action: string
  payload?: Record<string, unknown>
}
```

## 辅助函数

### Go

```go
// 创建字符串指针
func ptr[T any](v T) *T { return &v }

// 检查是否为标准组件类型
func IsStandardComponentType(typeName string) bool

// 获取组件类型名称
func (c *ComponentSpec) GetTypeName() string
```

### TypeScript

```typescript
// 创建消息处理器
function createMessageProcessor(registry?: ComponentRegistry): MessageProcessor

// 创建标准组件注册表
function createStandardRegistry(): ComponentRegistry

// 获取组件类型名称
function getComponentTypeName(spec: ComponentSpec): string

// 获取组件属性
function getComponentProps<T>(spec: ComponentSpec): T

// 检查是否为路径引用
function isPathReference(value: PropertyValue): boolean

// 文本清理（XSS 防护）
function sanitizeText(text: string): string

// URL 验证
function validateUrl(url: string): boolean
```

## 错误码

```typescript
const ErrorCodes = {
  INVALID_MESSAGE: 'INVALID_MESSAGE',      // 无效消息格式
  UNKNOWN_COMPONENT: 'UNKNOWN_COMPONENT',  // 未知组件类型
  INVALID_PATH: 'INVALID_PATH',            // 无效路径
  CIRCULAR_REFERENCE: 'CIRCULAR_REFERENCE', // 循环引用
  REGISTRY_FROZEN: 'REGISTRY_FROZEN',      // 注册表已冻结
  INVALID_TYPE_NAME: 'INVALID_TYPE_NAME',  // 无效类型名称
  XSS_DETECTED: 'XSS_DETECTED',            // 检测到 XSS
  INVALID_URL: 'INVALID_URL',              // 无效 URL
}
```
