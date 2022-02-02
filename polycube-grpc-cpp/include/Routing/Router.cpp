#include "Router.h"

#include "Details/PathTemplate.h"
#include "Details/RegexConverter.h"
#include "Exceptions.h"
#include <iostream>

namespace Routing {

using Details::PathTemplate;
using Details::RegexConverter;

Router::Router() {
}

void Router::registerPath(std::string const &path, std::string const &http_verb) {
    std::string uppercase_http_verb(http_verb);
    for (auto & c: uppercase_http_verb) c = std::toupper(c);
    _templates.push_back(PathTemplate(path,uppercase_http_verb));
}

PathMatch Router::matchPath(std::string const &path, std::string const &http_verb) {
    std::string uppercase_http_verb(http_verb);
    for (auto & c: uppercase_http_verb) c = std::toupper(c);
    for (auto &tpl: _templates) {
        if (std::regex_match(path, tpl.regex()) && tpl.http_verb() == uppercase_http_verb) {
            std::cout << "Inside matchPath: uppercase_http_verb " << uppercase_http_verb << " and path " << path << std::endl;
            return PathMatch(path, tpl,uppercase_http_verb);
        }
    }

    throw PathNotFoundException("Path not found for '" + path + "'");
}

}