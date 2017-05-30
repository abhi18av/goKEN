
Dir.chdir(File.dirname(__FILE__))
jsonFiles = Dir["*.json"]

for f in jsonFiles

newName1 =  f.split(".")[0] +  ".yaml" 

File.open(newName1, 'w')
cmd1 = "json2yaml "  + f + " > " + newName1
system(cmd1)

end
