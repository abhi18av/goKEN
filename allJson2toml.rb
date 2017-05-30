
Dir.chdir(File.dirname(__FILE__))
jsonFiles = Dir["*.json"]

for f in jsonFiles

newName1 =  f.split(".")[0] +  ".toml" 

#File.open(newName1, 'w')
cmd1 = "json2toml "  + f + " > " + newName1
system(cmd1)

end
