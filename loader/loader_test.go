package loader

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {

    l := New()

    TaskString := func (l *Loader, s string) Vertex {
        task := Task{
            Description: "creating string",
            Load: func(inputs Inputs) (any, error) {
                inputs.AssertNone()
                return s, nil
            },
        }
        return l.Add(task, nil, nil)
    }

    TaskJoin := func (l *Loader, separator Vertex, inputs ... Vertex) Vertex {
        task := Task{
            Description: "concatenating strings",
            Load: func(inputs Inputs) (any, error) {
                inputs.AssertNoErrors()

                sep := inputs.GetNamed("separator")

                var sb strings.Builder
                for i, v := range inputs.Direct {
                    sb.WriteString(v.Value.(string))
                    if i + 1 < len(inputs.Direct) {
                        sb.WriteString(sep.Value.(string))
                    }
                }

                return sb.String(), nil
            },
        }
        return l.Add(task, map[string]Vertex{"separator": separator}, inputs)
    }

    hello := TaskString(l, "hello")

    helloWorld := TaskJoin(l, // separator, inputs ...
        TaskString(l, ", "),
        hello,
        TaskString(l, "beautiful"),
        TaskString(l, "world!"),
    )

    hello.Keep()
    helloWorld.Keep()

    /*
    l.WalkSynchronous(func(t Task, i int, total int) {
        fmt.Printf("(%d/%d) %s...\n", i + 1, total, t.Description)
    })
     */

    l.Walk()

    // TODO
    // l.Result(helloWorld)
    assert.Equal(t, l.Result(helloWorld), Result{"hello, beautiful, world!", nil})

    // load...

    // foo := b.Result()

    /*

func LinkShaderTask(l *Loader, shaders ... *Vertex) *Vertex {
    t := loader.AddTask(Task{
        Load: func(inputs) Result {
            linkshader(input[i])
            linkshaders()
        }
    }, nil, shaders...)
}


gs_passthrough := CompileShaderTask(LoadFromPackTask("shaders.pack", "passthrough.gs"))

fs_skin := CompileShaderTask(TemplateTask(args, LoadFromPackTask("shaders.pack", "skin.fs")))
vs_skin := CompileShaderTask(TemplateTask(args, LoadFromPackTask("shaders.pack", "skin.vs")))

fs_fur := CompileShaderTask(TemplateTask(args, LoadFromPackTask("shaders.pack", "fur.fs")))
vs_fur := CompileShaderTask(TemplateTask(args, LoadFromPackTask("shaders.pack", "fur.vs")))

require("shader_skin", LinkShaderTask(fs_skin, vs_skin, gs_passthrough))
require("shader_fur",  LinkShaderTask(fs_fur,  vs_fur,  gs_passthrough))
keep("optional_bonus", DecodeImage(Decrypt(key, LoadFromPackTask("dlc1.pack", "optional_bonus.png"))))

require("foobar", FoobarTask(map[string]task{"foo": LoadFooTask, "bar": LoadBarTask}))

    ...

    later:

    shader_skin := GetRequiredResult("shader_skin")
    shader_fur  := GetRequiredResult("shader_fur")
    value, err := GetResult("optional_bonus")

     */

    /*

    DAG := digraph.New[Task, string, digraph.WeightDontCare]()

    a := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Load from disk", nil}
        },
    })
    b := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Decode from PNG", nil}
        },
    })
    c := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Upload as texture", nil}
        },
    })

    DAG.AddEdge(a, b, "")
    DAG.AddEdge(b, c, "")

    xf := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Load fragment shader from disk", nil}
        },
    })
    xv := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Load vertex shader from disk", nil}
        },
    })
    xg := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Load geometry shader from disk", nil}
        },
    })

    tf := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Template fragment shader", nil}
        },
    })
    tv := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Template vertex shader", nil}
        },
    })
    tg := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Template geometry shader", nil}
        },
    })

    yf := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Compile fragment shader", nil}
        },
    })
    yv := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Compile vertex shader", nil}
        },
    })
    yg := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Compile geometry shader", nil}
        },
    })
    z := DAG.AddVertex(Task{
        Load: func(inputs Inputs) *Result {
            return &Result{"Link shaders", nil}
        },
    })

    DAG.AddEdge(xf, tf, "")
    DAG.AddEdge(xv, tv, "")
    DAG.AddEdge(xg, tg, "")
    DAG.AddEdge(tf, yf, "")
    DAG.AddEdge(tv, yv, "")
    DAG.AddEdge(tg, yg, "")
    DAG.AddEdge(yf,  z, "")
    DAG.AddEdge(yv,  z, "")
    DAG.AddEdge(yg,  z, "")
     */
}
