import ballerina/io;
import celleryio/cellery;

public function build(cellery:ImageName iName) returns error? {
    //Build Stock Cell
    io:println("Building Stock Cell ...");
    //Stock Component
    cellery:Component stockComponent = {
        name: "stock",
        src: {
            image: "wso2cellery/sampleapp-stock:0.3.0"
        },
        ingresses: {
            stock: <cellery:HttpApiIngress>{ port: 8080,
                context: "stock",
                definition: {
                    resources: [
                        {
                            path: "/options",
                            method: "GET"
                        }
                    ]
                },
                expose: "local"
            }
        }
    };

    cellery:CellImage stockCell = {
        components: {
            stockComp: stockComponent
        }
    };
    return <@untainted> cellery:createImage(stockCell, iName);
}

public function run(cellery:ImageName iName, map<cellery:ImageName> instances, boolean startDependencies, boolean shareDependencies) returns (cellery:InstanceState[]|error?) {
    cellery:CellImage|cellery:Composite stockCell = cellery:constructImage(iName);
    return <@untainted> cellery:createInstance(stockCell, iName, instances, startDependencies, shareDependencies);
}
